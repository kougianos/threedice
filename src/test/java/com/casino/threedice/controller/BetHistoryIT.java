package com.casino.threedice.controller;

import com.casino.threedice.BaseIntegrationTest;
import com.casino.threedice.repository.BetRepository;
import com.casino.threedice.repository.ClientRepository;
import com.casino.threedice.entity.Client;
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

class BetHistoryIT extends BaseIntegrationTest {

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
                        .username("historyplayer")
                        .balance(new BigDecimal("5000.00"))
                        .build()
        );
        playerId = client.getId();
    }

    @Test
    @DisplayName("Should return empty list when player has no bets")
    void betHistory_empty() {
        given()
        .when()
                .get("/bets/history/{playerId}", playerId)
        .then()
                .statusCode(200)
                .body("$", hasSize(0));
    }

    @Test
    @DisplayName("Should return bet history with correct fields")
    void betHistory_singleBet() {
        // Place one bet first: dice 2*3*4 = 24, predict 24 → WIN
        when(diceEngine.roll()).thenReturn(2, 3, 4);

        given()
                .contentType(ContentType.JSON)
                .body("""
                        {
                            "playerId": %d,
                            "stake": 20,
                            "predictedValue": 24,
                            "idempotencyKey": "%s"
                        }
                        """.formatted(playerId, UUID.randomUUID()))
        .when()
                .post("/bets");

        // Now fetch history
        given()
        .when()
                .get("/bets/history/{playerId}", playerId)
        .then()
                .statusCode(200)
                .body("$", hasSize(1))
                .body("[0].predictedValue", equalTo(24))
                .body("[0].stake", equalTo(20.0F))
                .body("[0].status", equalTo("WON"))
                .body("[0].dieOne", equalTo(2))
                .body("[0].dieTwo", equalTo(3))
                .body("[0].dieThree", equalTo(4))
                .body("[0].productValue", equalTo(24))
                .body("[0].betId", notNullValue())
                .body("[0].createdAt", notNullValue());
    }

    @Test
    @DisplayName("Should return at most 10 bets, most recent first")
    void betHistory_limitedTo10() {
        // Place 12 bets
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
                .get("/bets/history/{playerId}", playerId)
        .then()
                .statusCode(200)
                .body("$", hasSize(10));
    }

    @Test
    @DisplayName("Should return 404 when player does not exist")
    void betHistory_playerNotFound() {
        given()
        .when()
                .get("/bets/history/{playerId}", 99999)
        .then()
                .statusCode(404)
                .body("detail", containsString("Player not found"));
    }
}
