package com.casino.threedice.dto;

import com.casino.threedice.entity.BetStatus;
import lombok.Builder;

import java.math.BigDecimal;

/**
 * Response payload after placing a bet.
 *
 * @param betId         the created bet ID
 * @param dieOne        first die value
 * @param dieTwo        second die value
 * @param dieThree      third die value
 * @param productValue  product of the three dice
 * @param status        WON or LOST
 * @param winnings      amount won (0 if lost)
 * @param balanceAfter  player's balance after the bet
 */
@Builder
public record PlaceBetResponse(
        Long betId,
        Integer dieOne,
        Integer dieTwo,
        Integer dieThree,
        Integer productValue,
        BetStatus status,
        BigDecimal winnings,
        BigDecimal balanceAfter
) {
}
