package com.casino.threedice.controller;

import com.casino.threedice.dto.TransactionHistoryEntry;
import com.casino.threedice.service.TransactionService;
import jakarta.validation.constraints.Positive;
import lombok.RequiredArgsConstructor;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@Validated
@RestController
@RequestMapping("/api/transactions")
@RequiredArgsConstructor
public class TransactionController {

    private final TransactionService transactionService;

    /**
     * Returns the last 10 transactions for the given player, most recent first.
     */
    @GetMapping("/history/{playerId}")
    public List<TransactionHistoryEntry> getTransactionHistory(
            @PathVariable @Positive(message = "Player ID must be positive") Long playerId) {
        return transactionService.getTransactionHistory(playerId);
    }
}
