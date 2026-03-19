package com.casino.threedice.controller;

import com.casino.threedice.BaseIntegrationTest;
import com.casino.threedice.entity.Client;
import com.casino.threedice.repository.BetRepository;
import com.casino.threedice.repository.ClientRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;

import java.math.BigDecimal;

import static io.restassured.RestAssured.given;
import static org.hamcrest.Matchers.*;

class PlayerControllerIT extends BaseIntegrationTest {

    @Autowired
    private ClientRepository clientRepository;

    @Autowired
    private BetRepository betRepository;

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

    @Test
    @DisplayName("Should return player details with current balance")
    void getPlayer_success() {
        given()
        .when()
                .get("/players/{playerId}", playerId)
        .then()
                .statusCode(200)
                .body("id", equalTo(playerId.intValue()))
                .body("username", equalTo("testplayer"))
                .body("balance", equalTo(1000.0F))
                .body("createdAt", notNullValue());
    }

    @Test
    @DisplayName("Should return 404 for non-existent player")
    void getPlayer_notFound() {
        given()
        .when()
                .get("/players/{playerId}", 99999)
        .then()
                .statusCode(404)
                .body("detail", containsString("Player not found"));
    }
}
