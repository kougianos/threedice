package com.casino.threedice.exception;

public class DuplicateBetException extends RuntimeException {

    public DuplicateBetException(String idempotencyKey) {
        super("Duplicate bet submission: idempotencyKey=%s".formatted(idempotencyKey));
    }
}
