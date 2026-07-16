package com.casino.threedice.exception;

import java.math.BigDecimal;

public class InsufficientBalanceException extends RuntimeException {

    public InsufficientBalanceException(BigDecimal balance, BigDecimal stake) {
        super("Insufficient balance: current=%s, requested stake=%s".formatted(balance, stake));
    }
}
