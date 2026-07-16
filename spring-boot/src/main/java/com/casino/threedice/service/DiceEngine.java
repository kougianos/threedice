package com.casino.threedice.service;

import org.springframework.stereotype.Component;

import java.math.BigDecimal;
import java.security.SecureRandom;
import java.util.Collections;
import java.util.Set;
import java.util.TreeSet;

/**
 * Encapsulates all dice-related logic: rolling, valid product computation, and odds calculation.
 *
 * <p>Valid products: only 40 out of 216 combinations are achievable with three six-sided dice.
 *
 * <p>Odds:
 *   Product < 9         -> x2
 *   9 ≤ Product < 120   -> x5
 *   Product ≥ 120       -> x2
 */
@Component
public class DiceEngine {

    private static final SecureRandom SECURE_RANDOM = new SecureRandom();
    private static final int DICE_SIDES = 6;

    public static final Set<Integer> VALID_PRODUCTS;

    static {
        TreeSet<Integer> products = new TreeSet<>();
        for (int a = 1; a <= 6; a++) {
            for (int b = 1; b <= 6; b++) {
                for (int c = 1; c <= 6; c++) {
                    products.add(a * b * c);
                }
            }
        }
        VALID_PRODUCTS = Collections.unmodifiableSet(products);
    }

    /**
     * Rolls a single six-sided die using a cryptographically secure random source.
     *
     * @return a value between 1 and 6 inclusive
     */
    public int roll() {
        return SECURE_RANDOM.nextInt(DICE_SIDES) + 1;
    }

    /**
     * Returns whether the given value is mathematically achievable as a product of three dice.
     */
    public boolean isValidProduct(int value) {
        return VALID_PRODUCTS.contains(value);
    }

    /**
     * Returns the odds multiplier for a given product value.
     */
    public static BigDecimal getOdds(int productValue) {
        if (productValue < 9) {
            return BigDecimal.valueOf(2);
        } else if (productValue < 120) {
            return BigDecimal.valueOf(5);
        } else {
            return BigDecimal.valueOf(2);
        }
    }
}
