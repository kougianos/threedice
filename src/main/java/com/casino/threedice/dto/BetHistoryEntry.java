package com.casino.threedice.dto;

import com.casino.threedice.entity.BetStatus;
import lombok.Builder;

import java.math.BigDecimal;
import java.time.Instant;

@Builder
public record BetHistoryEntry(
        Long betId,
        Integer predictedValue,
        BigDecimal stake,
        BetStatus status,
        Integer dieOne,
        Integer dieTwo,
        Integer dieThree,
        Integer productValue,
        Instant createdAt
) {
}
