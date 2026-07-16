package com.casino.threedice.entity;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.math.BigDecimal;
import java.time.Instant;

/**
 * Represents a single bet placed by a client.
 * Links to the draw result and the financial transaction.
 */
@Entity
@Table(name = "bets")
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class Bet {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "client_id", nullable = false)
    private Client client;

    @Column(name = "predicted_value", nullable = false)
    private Integer predictedValue;

    @Column(nullable = false, precision = 12, scale = 2)
    private BigDecimal stake;

    @Column(name = "idempotency_key", nullable = false, unique = true, length = 64)
    private String idempotencyKey;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false, length = 4)
    private BetStatus status;

    @OneToOne(mappedBy = "bet", cascade = CascadeType.ALL)
    private Draw draw;

    @OneToOne(mappedBy = "bet", cascade = CascadeType.ALL)
    private Transaction transaction;

    @Column(name = "created_at", nullable = false, updatable = false)
    private Instant createdAt;

    @PrePersist
    void prePersist() {
        this.createdAt = Instant.now();
    }
}
