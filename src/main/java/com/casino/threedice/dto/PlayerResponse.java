package com.casino.threedice.dto;

import lombok.Builder;

import java.math.BigDecimal;
import java.time.Instant;

@Builder
public record PlayerResponse(
        Long id,
        String username,
        BigDecimal balance,
        Instant createdAt
) {
}
