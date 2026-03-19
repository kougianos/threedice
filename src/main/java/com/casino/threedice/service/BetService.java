package com.casino.threedice.service;

import com.casino.threedice.dto.BetHistoryEntry;
import com.casino.threedice.dto.PlaceBetRequest;
import com.casino.threedice.dto.PlaceBetResponse;
import com.casino.threedice.entity.*;
import com.casino.threedice.exception.DuplicateBetException;
import com.casino.threedice.exception.InsufficientBalanceException;
import com.casino.threedice.exception.PlayerNotFoundException;
import com.casino.threedice.repository.BetRepository;
import com.casino.threedice.repository.ClientRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.PageRequest;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.util.List;

@Slf4j
@Service
@RequiredArgsConstructor
public class BetService {

    private final ClientRepository clientRepository;
    private final BetRepository betRepository;
    private final DiceEngine diceEngine;

    /**
     * Places a bet for the given player: deducts the stake, rolls dice,
     * determines win/loss, updates balance, and persists all records atomically.
     */
    @Transactional
    public PlaceBetResponse placeBet(PlaceBetRequest request) {
        // Reject duplicate submissions early
        if (betRepository.existsByIdempotencyKey(request.idempotencyKey())) {
            throw new DuplicateBetException(request.idempotencyKey());
        }

        // Pessimistic lock: SELECT ... FOR UPDATE prevents concurrent balance reads
        Client client = clientRepository.findByIdForUpdate(request.playerId())
                .orElseThrow(() -> new PlayerNotFoundException(request.playerId()));

        // Validate sufficient balance
        if (client.getBalance().compareTo(request.stake()) < 0) {
            throw new InsufficientBalanceException(client.getBalance(), request.stake());
        }

        // Deduct stake immediately
        client.setBalance(client.getBalance().subtract(request.stake()));

        // Roll three dice
        int d1 = diceEngine.roll();
        int d2 = diceEngine.roll();
        int d3 = diceEngine.roll();
        int productValue = d1 * d2 * d3;

        // Determine outcome
        boolean won = productValue == request.predictedValue();
        BetStatus status = won ? BetStatus.WON : BetStatus.LOST;

        // Calculate winnings and update balance if won
        BigDecimal winnings = BigDecimal.ZERO;
        if (won) {
            BigDecimal odds = DiceEngine.getOdds(productValue);
            winnings = request.stake().multiply(odds);
            client.setBalance(client.getBalance().add(winnings));
        }

        clientRepository.save(client);

        // Create and persist bet
        Bet bet = Bet.builder()
                .client(client)
                .predictedValue(request.predictedValue())
                .stake(request.stake())
                .idempotencyKey(request.idempotencyKey())
                .status(status)
                .build();

        // Create draw linked to bet
        Draw draw = Draw.builder()
                .bet(bet)
                .dieOne(d1)
                .dieTwo(d2)
                .dieThree(d3)
                .productValue(productValue)
                .build();
        bet.setDraw(draw);

        // Create transaction record.
        // On a loss the transaction is a DEBIT of the stake.
        // On a win we record a CREDIT of the net winnings (winnings - stake already deducted).
        Transaction transaction = Transaction.builder()
                .bet(bet)
                .client(client)
                .type(won ? TransactionType.CREDIT : TransactionType.DEBIT)
                .amount(won ? winnings : request.stake())
                .balanceAfter(client.getBalance())
                .build();
        bet.setTransaction(transaction);

        betRepository.save(bet);

        log.info("Bet placed: playerId={}, predicted={}, actual={}, status={}, balance={}",
                client.getId(), request.predictedValue(), productValue, status, client.getBalance());

        return PlaceBetResponse.builder()
                .betId(bet.getId())
                .dieOne(d1)
                .dieTwo(d2)
                .dieThree(d3)
                .productValue(productValue)
                .status(status)
                .winnings(winnings)
                .balanceAfter(client.getBalance())
                .build();
    }

    /**
     * Returns the last 10 bets for a given player, most recent first.
     */
    @Transactional(readOnly = true)
    public List<BetHistoryEntry> getBetHistory(Long playerId) {
        validatePlayerExists(playerId);

        return betRepository.findByClientIdOrderByCreatedAtDesc(playerId, PageRequest.of(0, 10))
                .stream()
                .map(bet -> BetHistoryEntry.builder()
                        .betId(bet.getId())
                        .predictedValue(bet.getPredictedValue())
                        .stake(bet.getStake())
                        .status(bet.getStatus())
                        .dieOne(bet.getDraw().getDieOne())
                        .dieTwo(bet.getDraw().getDieTwo())
                        .dieThree(bet.getDraw().getDieThree())
                        .productValue(bet.getDraw().getProductValue())
                        .createdAt(bet.getCreatedAt())
                        .build())
                .toList();
    }

    private void validatePlayerExists(Long playerId) {
        if (!clientRepository.existsById(playerId)) {
            throw new PlayerNotFoundException(playerId);
        }
    }
}
