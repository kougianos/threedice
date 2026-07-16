package com.casino.threedice.dto;

import com.casino.threedice.entity.TransactionType;
import lombok.Builder;

import java.math.BigDecimal;
import java.time.Instant;

@Builder
public record TransactionHistoryEntry(
        Long transactionId,
        Long betId,
        TransactionType type,
        BigDecimal amount,
        BigDecimal balanceAfter,
        Instant createdAt
) {
}
