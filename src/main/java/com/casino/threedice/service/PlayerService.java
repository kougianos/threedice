package com.casino.threedice.service;

import com.casino.threedice.dto.PlayerResponse;
import com.casino.threedice.entity.Client;
import com.casino.threedice.exception.PlayerNotFoundException;
import com.casino.threedice.repository.ClientRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.CommandLineRunner;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;

@Slf4j
@Service
@RequiredArgsConstructor
public class PlayerService implements CommandLineRunner {

    private final ClientRepository clientRepository;

    @Value("${game.initial-balance}")
    private BigDecimal initialBalance;

    /**
     * Seeds a default demo player on startup if no clients exist.
     */
    @Override
    public void run(String... args) {
        if (clientRepository.count() == 0) {
            Client demo = clientRepository.save(
                    Client.builder()
                            .username("player1")
                            .balance(initialBalance)
                            .build()
            );
            log.info("Created demo player: id={}, username={}, balance={}",
                    demo.getId(), demo.getUsername(), demo.getBalance());
        }
    }

    @Transactional(readOnly = true)
    public PlayerResponse getPlayer(Long playerId) {
        return clientRepository.findById(playerId)
                .map(c -> PlayerResponse.builder()
                        .id(c.getId())
                        .username(c.getUsername())
                        .balance(c.getBalance())
                        .createdAt(c.getCreatedAt())
                        .build())
                .orElseThrow(() -> new PlayerNotFoundException(playerId));
    }
}
