package com.casino.threedice.controller;

import com.casino.threedice.dto.PlayerResponse;
import com.casino.threedice.service.PlayerService;
import jakarta.validation.constraints.Positive;
import lombok.RequiredArgsConstructor;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

@Validated
@RestController
@RequestMapping("/api/players")
@RequiredArgsConstructor
public class PlayerController {

    private final PlayerService playerService;

    /**
     * Returns the player's current information including balance.
     */
    @GetMapping("/{playerId}")
    public PlayerResponse getPlayer(
            @PathVariable @Positive(message = "Player ID must be positive") Long playerId) {
        return playerService.getPlayer(playerId);
    }
}
