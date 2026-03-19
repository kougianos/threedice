package com.casino.threedice.controller;

import com.casino.threedice.BaseIntegrationTest;
import com.casino.threedice.entity.Client;
import com.casino.threedice.repository.BetRepository;
import com.casino.threedice.repository.ClientRepository;
import com.casino.threedice.service.DiceEngine;
import io.restassured.http.ContentType;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.test.context.bean.override.mockito.MockitoBean;

import java.math.BigDecimal;
import java.util.UUID;

import static io.restassured.RestAssured.given;
import static org.hamcrest.Matchers.*;
import static org.mockito.Mockito.when;

class TransactionHistoryIT extends BaseIntegrationTest {

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
                        .username("txnplayer")
                        .balance(new BigDecimal("5000.00"))
                        .build()
        );
        playerId = client.getId();
    }

    @Test
    @DisplayName("Should return empty list when player has no transactions")
    void transactionHistory_empty() {
        given()
        .when()
                .get("/transactions/history/{playerId}", playerId)
        .then()
                .statusCode(200)
                .body("$", hasSize(0));
    }

    @Test
    @DisplayName("Should return DEBIT transaction on a lost bet")
    void transactionHistory_debitOnLoss() {
        // Dice 5*5*5 = 125, predicted 12 → LOST → DEBIT
        when(diceEngine.roll()).thenReturn(5, 5, 5);

        given()
                .contentType(ContentType.JSON)
                .body("""
                        {
                            "playerId": %d,
                            "stake": 30,
                            "predictedValue": 12,
                            "idempotencyKey": "%s"
                        }
                        """.formatted(playerId, UUID.randomUUID()))
        .when()
                .post("/bets");

        given()
        .when()
                .get("/transactions/history/{playerId}", playerId)
        .then()
                .statusCode(200)
                .body("$", hasSize(1))
                .body("[0].type", equalTo("DEBIT"))
                .body("[0].amount", equalTo(30.0F))
                .body("[0].balanceAfter", equalTo(4970.0F))
                .body("[0].betId", notNullValue())
                .body("[0].transactionId", notNullValue())
                .body("[0].createdAt", notNullValue());
    }

    @Test
    @DisplayName("Should return CREDIT transaction on a won bet")
    void transactionHistory_creditOnWin() {
        // Dice 2*3*4 = 24, predicted 24 → WON, odds 5, winnings = 50*5 = 250
        // Balance: 5000 - 50 + 250 = 5200
        when(diceEngine.roll()).thenReturn(2, 3, 4);

        given()
                .contentType(ContentType.JSON)
                .body("""
                        {
                            "playerId": %d,
                            "stake": 50,
                            "predictedValue": 24,
                            "idempotencyKey": "%s"
                        }
                        """.formatted(playerId, UUID.randomUUID()))
        .when()
                .post("/bets");

        given()
        .when()
                .get("/transactions/history/{playerId}", playerId)
        .then()
                .statusCode(200)
                .body("$", hasSize(1))
                .body("[0].type", equalTo("CREDIT"))
                .body("[0].amount", equalTo(250.0F))
                .body("[0].balanceAfter", equalTo(5200.0F));
    }

    @Test
    @DisplayName("Should return at most 10 transactions, most recent first")
    void transactionHistory_limitedTo10() {
        when(diceEngine.roll()).thenReturn(1, 1, 1); // product = 1

        for (int i = 0; i < 12; i++) {
            given()
                    .contentType(ContentType.JSON)
                    .body("""
                            {
                                "playerId": %d,
                                "stake": 5,
                                "predictedValue": 100,
                                "idempotencyKey": "%s"
                            }
                            """.formatted(playerId, UUID.randomUUID()))
            .when()
                    .post("/bets");
        }

        given()
        .when()
                .get("/transactions/history/{playerId}", playerId)
        .then()
                .statusCode(200)
                .body("$", hasSize(10));
    }

    @Test
    @DisplayName("Should return 404 when player does not exist")
    void transactionHistory_playerNotFound() {
        given()
        .when()
                .get("/transactions/history/{playerId}", 99999)
        .then()
                .statusCode(404)
                .body("detail", containsString("Player not found"));
    }
}
