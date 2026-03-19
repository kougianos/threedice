package com.casino.threedice.exception;

public class PlayerNotFoundException extends RuntimeException {

    public PlayerNotFoundException(Long playerId) {
        super("Player not found with ID: " + playerId);
    }
}
