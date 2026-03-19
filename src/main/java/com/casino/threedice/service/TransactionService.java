package com.casino.threedice.service;

import com.casino.threedice.dto.TransactionHistoryEntry;
import com.casino.threedice.exception.PlayerNotFoundException;
import com.casino.threedice.repository.ClientRepository;
import com.casino.threedice.repository.TransactionRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.PageRequest;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Slf4j
@Service
@RequiredArgsConstructor
public class TransactionService {

    private final TransactionRepository transactionRepository;
    private final ClientRepository clientRepository;

    /**
     * Returns the last 10 transactions for a given player, most recent first.
     */
    @Transactional(readOnly = true)
    public List<TransactionHistoryEntry> getTransactionHistory(Long playerId) {
        if (!clientRepository.existsById(playerId)) {
            throw new PlayerNotFoundException(playerId);
        }

        return transactionRepository.findByClientIdOrderByCreatedAtDesc(playerId, PageRequest.of(0, 10))
                .stream()
                .map(tx -> TransactionHistoryEntry.builder()
                        .transactionId(tx.getId())
                        .betId(tx.getBet().getId())
                        .type(tx.getType())
                        .amount(tx.getAmount())
                        .balanceAfter(tx.getBalanceAfter())
                        .createdAt(tx.getCreatedAt())
                        .build())
                .toList();
    }
}
