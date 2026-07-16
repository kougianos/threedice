package com.casino.threedice.validation;

import jakarta.validation.Constraint;
import jakarta.validation.Payload;

import java.lang.annotation.*;

/**
 * Validates that a predicted product value is mathematically achievable
 * by rolling three six-sided dice.
 */
@Documented
@Constraint(validatedBy = ValidPredictedValueValidator.class)
@Target({ElementType.FIELD, ElementType.PARAMETER})
@Retention(RetentionPolicy.RUNTIME)
public @interface ValidPredictedValue {

    String message() default "Predicted value is not a possible product of three dice";

    Class<?>[] groups() default {};

    Class<? extends Payload>[] payload() default {};
}
