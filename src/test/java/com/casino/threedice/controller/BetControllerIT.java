package com.casino.threedice.controller;

import com.casino.threedice.BaseIntegrationTest;
import com.casino.threedice.entity.Client;
import com.casino.threedice.repository.BetRepository;
import com.casino.threedice.repository.ClientRepository;
import com.casino.threedice.service.DiceEngine;
import io.restassured.http.ContentType;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.test.context.bean.override.mockito.MockitoBean;

import java.math.BigDecimal;
import java.util.UUID;

import static io.restassured.RestAssured.given;
import static org.hamcrest.Matchers.*;
import static org.mockito.Mockito.when;

class BetControllerIT extends BaseIntegrationTest {

    @Autowired
    private ClientRepository clientRepository;

    @Autowired
    private BetRepository betRepository;


    @MockitoBean
    private DiceEngine diceEngine;

    private Long playerId;

    @BeforeEach
    void setUp() {
        betRepository.deleteAll();
        clientRepository.deleteAll();

        Client client = clientRepository.save(
                Client.builder()
                        .username("testplayer")
                        .balance(new BigDecimal("1000.00"))
                        .build()
        );
        playerId = client.getId();
    }

    // ───── Core Bet Placement ─────

    @Test
    @DisplayName("Should win bet when predicted value matches dice product")
    void placeBet_win() {
        when(diceEngine.roll()).thenReturn(1, 3, 4);

        given()
                .contentType(ContentType.JSON)
                .body("""
                        {
                            "playerId": %d,
                            "stake": 10,
                            "predictedValue": 12,
                            "idempotencyKey": "%s"
                        }
                        """.formatted(playerId, UUID.randomUUID()))
        .when()
                .post("/bets")
        .then()
                .statusCode(201)
                .body("dieOne", equalTo(1))
                .body("dieTwo", equalTo(3))
                .body("dieThree", equalTo(4))
                .body("productValue", equalTo(12))
                .body("status", equalTo("WON"))
                .body("winnings", equalTo(50))
                .body("balanceAfter", equalTo(1040.0F));
    }

    @Test
    @DisplayName("Should lose bet when predicted value does not match dice product")
    void placeBet_lose() {
        when(diceEngine.roll()).thenReturn(2, 3, 5);

        given()
                .contentType(ContentType.JSON)
                .body("""
                        {
                            "playerId": %d,
                            "stake": 10,
                            "predictedValue": 12,
                            "idempotencyKey": "%s"
                        }
                        """.formatted(playerId, UUID.randomUUID()))
        .when()
                .post("/bets")
        .then()
                .statusCode(201)
                .body("status", equalTo("LOST"))
                .body("productValue", equalTo(30))
                .body("winnings", equalTo(0))
                .body("balanceAfter", equalTo(990.0F));
    }

    @Test
    @DisplayName("Balance should update correctly across multiple bets")
    void placeBet_multipleRounds() {
        when(diceEngine.roll()).thenReturn(6, 6, 6);

        given()
                .contentType(ContentType.JSON)
                .body("""
                        {
                            "playerId": %d,
                            "stake": 50,
                            "predictedValue": 12,
                            "idempotencyKey": "%s"
                        }
                        """.formatted(playerId, UUID.randomUUID()))
        .when()
                .post("/bets")
        .then()
                .statusCode(201)
                .body("status", equalTo("LOST"))
                .body("balanceAfter", equalTo(950.0F));

        when(diceEngine.roll()).thenReturn(2, 2, 2);

        given()
                .contentType(ContentType.JSON)
                .body("""
                        {
                            "playerId": %d,
                            "stake": 100,
                            "predictedValue": 8,
                            "idempotencyKey": "%s"
                        }
                        """.formatted(playerId, UUID.randomUUID()))
        .when()
                .post("/bets")
        .then()
                .statusCode(201)
                .body("status", equalTo("WON"))
                .body("winnings", equalTo(200))
                .body("balanceAfter", equalTo(1050.0F));
    }

    // ───── Validation ─────

    @Nested
    @DisplayName("Request Validation")
    class ValidationTests {

        @Test
        @DisplayName("Should return 404 when player does not exist")
        void placeBet_playerNotFound() {
            given()
                    .contentType(ContentType.JSON)
                    .body("""
                            {
                                "playerId": 99999,
                                "stake": 10,
                                "predictedValue": 12,
                                "idempotencyKey": "%s"
                            }
                            """.formatted(UUID.randomUUID()))
            .when()
                    .post("/bets")
            .then()
                    .statusCode(404)
                    .body("detail", containsString("Player not found"));
        }

        @Test
        @DisplayName("Should return 400 when stake exceeds balance")
        void placeBet_insufficientBalance() {
            given()
                    .contentType(ContentType.JSON)
                    .body("""
                            {
                                "playerId": %d,
                                "stake": 5000,
                                "predictedValue": 12,
                                "idempotencyKey": "%s"
                            }
                            """.formatted(playerId, UUID.randomUUID()))
            .when()
                    .post("/bets")
            .then()
                    .statusCode(400)
                    .body("detail", containsString("Insufficient balance"));
        }

        @Test
        @DisplayName("Should return 400 when request validation fails - missing fields")
        void placeBet_validationError_missingFields() {
            given()
                    .contentType(ContentType.JSON)
                    .body("{}")
            .when()
                    .post("/bets")
            .then()
                    .statusCode(400)
                    .body("detail", containsString("Player ID is required"));
        }

        @Test
        @DisplayName("Should return 400 when predicted value is out of range")
        void placeBet_validationError_predictedValueOutOfRange() {
            given()
                    .contentType(ContentType.JSON)
                    .body("""
                            {
                                "playerId": %d,
                                "stake": 10,
                                "predictedValue": 300,
                                "idempotencyKey": "%s"
                            }
                            """.formatted(playerId, UUID.randomUUID()))
            .when()
                    .post("/bets")
            .then()
                    .statusCode(400)
                    .body("detail", containsString("Predicted value must be at most 216"));
        }

        @Test
        @DisplayName("Should return 400 when stake is negative")
        void placeBet_validationError_negativeStake() {
            given()
                    .contentType(ContentType.JSON)
                    .body("""
                            {
                                "playerId": %d,
                                "stake": -5,
                                "predictedValue": 12,
                                "idempotencyKey": "%s"
                            }
                            """.formatted(playerId, UUID.randomUUID()))
            .when()
                    .post("/bets")
            .then()
                    .statusCode(400)
                    .body("detail", containsString("Minimum stake is $1.00"));
        }

        @Test
        @DisplayName("Should return 400 when predicted value is an impossible dice product (e.g. 7)")
        void placeBet_validationError_impossibleProduct() {
            given()
                    .contentType(ContentType.JSON)
                    .body("""
                            {
                                "playerId": %d,
                                "stake": 10,
                                "predictedValue": 7,
                                "idempotencyKey": "%s"
                            }
                            """.formatted(playerId, UUID.randomUUID()))
            .when()
                    .post("/bets")
            .then()
                    .statusCode(400)
                    .body("detail", containsString("not a possible product of three dice"));
        }

        @Test
        @DisplayName("Should accept valid edge-case products (1 and 216)")
        void placeBet_validEdgeCaseProducts() {
            when(diceEngine.roll()).thenReturn(1, 1, 1);

            given()
                    .contentType(ContentType.JSON)
                    .body("""
                            {
                                "playerId": %d,
                                "stake": 1,
                                "predictedValue": 1,
                                "idempotencyKey": "%s"
                            }
                            """.formatted(playerId, UUID.randomUUID()))
            .when()
                    .post("/bets")
            .then()
                    .statusCode(201);

            when(diceEngine.roll()).thenReturn(6, 6, 6);

            given()
                    .contentType(ContentType.JSON)
                    .body("""
                            {
                                "playerId": %d,
                                "stake": 1,
                                "predictedValue": 216,
                                "idempotencyKey": "%s"
                            }
                            """.formatted(playerId, UUID.randomUUID()))
            .when()
                    .post("/bets")
            .then()
                    .statusCode(201);
        }
    }

    // ───── Stake Limits ─────

    @Nested
    @DisplayName("Stake Limits")
    class StakeLimitTests {

        @Test
        @DisplayName("Should reject stake below minimum ($1.00)")
        void stakeBelowMinimum_shouldReturn400() {
            given()
                    .contentType(ContentType.JSON)
                    .body("""
                            {
                                "playerId": %d,
                                "stake": 0.50,
                                "predictedValue": 12,
                                "idempotencyKey": "%s"
                            }
                            """.formatted(playerId, UUID.randomUUID()))
            .when()
                    .post("/bets")
            .then()
                    .statusCode(400)
                    .body("detail", containsString("Minimum stake is $1.00"));
        }

        @Test
        @DisplayName("Should reject stake above maximum ($10,000.00)")
        void stakeAboveMaximum_shouldReturn400() {
            given()
                    .contentType(ContentType.JSON)
                    .body("""
                            {
                                "playerId": %d,
                                "stake": 10001,
                                "predictedValue": 12,
                                "idempotencyKey": "%s"
                            }
                            """.formatted(playerId, UUID.randomUUID()))
            .when()
                    .post("/bets")
            .then()
                    .statusCode(400)
                    .body("detail", containsString("Maximum stake is $10,000.00"));
        }

        @Test
        @DisplayName("Should accept stake at exact minimum ($1.00)")
        void stakeAtMinimum_shouldSucceed() {
            when(diceEngine.roll()).thenReturn(1, 1, 1);

            given()
                    .contentType(ContentType.JSON)
                    .body("""
                            {
                                "playerId": %d,
                                "stake": 1.00,
                                "predictedValue": 12,
                                "idempotencyKey": "%s"
                            }
                            """.formatted(playerId, UUID.randomUUID()))
            .when()
                    .post("/bets")
            .then()
                    .statusCode(201)
                    .body("balanceAfter", equalTo(999.0F));
        }

        @Test
        @DisplayName("Should reject max stake when it exceeds balance")
        void stakeAtMaximum_shouldBeRejectedByBalance() {
            when(diceEngine.roll()).thenReturn(1, 1, 1);

            given()
                    .contentType(ContentType.JSON)
                    .body("""
                            {
                                "playerId": %d,
                                "stake": 10000,
                                "predictedValue": 12,
                                "idempotencyKey": "%s"
                            }
                            """.formatted(playerId, UUID.randomUUID()))
            .when()
                    .post("/bets")
            .then()
                    .statusCode(400)
                    .body("detail", containsString("Insufficient balance"));
        }
    }

    // ───── Idempotency ─────

    @Nested
    @DisplayName("Idempotency Key")
    class IdempotencyTests {

        @Test
        @DisplayName("Should reject duplicate bet with same idempotency key")
        void duplicateBet_shouldReturn409() {
            when(diceEngine.roll()).thenReturn(1, 1, 1);
            String idempotencyKey = UUID.randomUUID().toString();

            given()
                    .contentType(ContentType.JSON)
                    .body("""
                            {
                                "playerId": %d,
                                "stake": 10,
                                "predictedValue": 12,
                                "idempotencyKey": "%s"
                            }
                            """.formatted(playerId, idempotencyKey))
            .when()
                    .post("/bets")
            .then()
                    .statusCode(201);

            given()
                    .contentType(ContentType.JSON)
                    .body("""
                            {
                                "playerId": %d,
                                "stake": 10,
                                "predictedValue": 12,
                                "idempotencyKey": "%s"
                            }
                            """.formatted(playerId, idempotencyKey))
            .when()
                    .post("/bets")
            .then()
                    .statusCode(409)
                    .body("detail", containsString("Duplicate bet submission"));
        }

        @Test
        @DisplayName("Should accept bets with different idempotency keys")
        void differentKeys_shouldBothSucceed() {
            when(diceEngine.roll()).thenReturn(1, 1, 1);

            given()
                    .contentType(ContentType.JSON)
                    .body("""
                            {
                                "playerId": %d,
                                "stake": 10,
                                "predictedValue": 12,
                                "idempotencyKey": "%s"
                            }
                            """.formatted(playerId, UUID.randomUUID()))
            .when()
                    .post("/bets")
            .then()
                    .statusCode(201);

            given()
                    .contentType(ContentType.JSON)
                    .body("""
                            {
                                "playerId": %d,
                                "stake": 10,
                                "predictedValue": 12,
                                "idempotencyKey": "%s"
                            }
                            """.formatted(playerId, UUID.randomUUID()))
            .when()
                    .post("/bets")
            .then()
                    .statusCode(201);
        }

        @Test
        @DisplayName("Should return 400 when idempotency key is missing")
        void missingIdempotencyKey_shouldReturn400() {
            given()
                    .contentType(ContentType.JSON)
                    .body("""
                            {
                                "playerId": %d,
                                "stake": 10,
                                "predictedValue": 12
                            }
                            """.formatted(playerId))
            .when()
                    .post("/bets")
            .then()
                    .statusCode(400)
                    .body("detail", containsString("Idempotency key is required"));
        }
    }

    // ───── Balance Consistency ─────

    @Test
    @DisplayName("Balance should remain consistent after rapid sequential bets")
    void rapidSequentialBets_balanceShouldBeConsistent() {
        when(diceEngine.roll()).thenReturn(6, 6, 6);

        for (int i = 0; i < 10; i++) {
            given()
                    .contentType(ContentType.JSON)
                    .body("""
                            {
                                "playerId": %d,
                                "stake": 50,
                                "predictedValue": 12,
                                "idempotencyKey": "%s"
                            }
                            """.formatted(playerId, UUID.randomUUID()))
            .when()
                    .post("/bets")
            .then()
                    .statusCode(201);
        }

        given()
        .when()
                .get("/players/{playerId}", playerId)
        .then()
                .statusCode(200)
                .body("balance", equalTo(500.0F));
    }
}
