package com.casino.threedice.controller;

import com.casino.threedice.dto.BetHistoryEntry;
import com.casino.threedice.dto.PlaceBetRequest;
import com.casino.threedice.dto.PlaceBetResponse;
import com.casino.threedice.service.BetService;
import jakarta.validation.Valid;
import jakarta.validation.constraints.Positive;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@Validated
@RestController
@RequestMapping("/api/bets")
@RequiredArgsConstructor
public class BetController {

    private final BetService betService;

    /**
     * Place a new bet. Rolls three dice, determines outcome,
     * updates the player's balance, and returns the result.
     */
    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public PlaceBetResponse placeBet(@Valid @RequestBody PlaceBetRequest request) {
        return betService.placeBet(request);
    }

    /**
     * Returns the last 10 bets for the given player, most recent first.
     */
    @GetMapping("/history/{playerId}")
    public List<BetHistoryEntry> getBetHistory(
            @PathVariable @Positive(message = "Player ID must be positive") Long playerId) {
        return betService.getBetHistory(playerId);
    }
}
