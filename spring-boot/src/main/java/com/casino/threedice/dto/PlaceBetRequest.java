package com.casino.threedice.dto;

import com.casino.threedice.validation.ValidPredictedValue;
import jakarta.validation.constraints.*;
import lombok.Builder;

import java.math.BigDecimal;

/**
 * Request payload for placing a bet.
 *
 * @param playerId       the client placing the bet
 * @param stake          the amount wagered (must be between 1 and 10,000)
 * @param predictedValue the predicted product of three dice (valid range: 1–216)
 * @param idempotencyKey unique key to prevent duplicate bet submissions
 */
@Builder
public record PlaceBetRequest(
        @NotNull(message = "Player ID is required")
        Long playerId,

        @NotNull(message = "Stake is required")
        @DecimalMin(value = "1.00", message = "Minimum stake is $1.00")
        @DecimalMax(value = "10000.00", message = "Maximum stake is $10,000.00")
        BigDecimal stake,

        @NotNull(message = "Predicted value is required")
        @Min(value = 1, message = "Predicted value must be at least 1")
        @Max(value = 216, message = "Predicted value must be at most 216")
        @ValidPredictedValue
        Integer predictedValue,

        @NotBlank(message = "Idempotency key is required")
        @Size(max = 64, message = "Idempotency key must be at most 64 characters")
        String idempotencyKey
) {
}
