package com.casino.threedice.validation;

import com.casino.threedice.service.DiceEngine;
import jakarta.validation.ConstraintValidator;
import jakarta.validation.ConstraintValidatorContext;

/**
 * Checks whether the predicted value is in the set of mathematically possible
 * products of three six-sided dice.
 */
public class ValidPredictedValueValidator implements ConstraintValidator<ValidPredictedValue, Integer> {

    @Override
    public boolean isValid(Integer value, ConstraintValidatorContext context) {
        // let @NotNull on the field handle null rejection
        if (value == null) {
            return true;
        }

        if (DiceEngine.VALID_PRODUCTS.contains(value)) {
            return true;
        }

        context.disableDefaultConstraintViolation();
        context.buildConstraintViolationWithTemplate(
            "Predicted value " + value + " is not a possible product of three dice"
        ).addConstraintViolation();
        return false;
    }
}
