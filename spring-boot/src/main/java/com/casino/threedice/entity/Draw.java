package com.casino.threedice.entity;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

/**
 * Stores the three dice values drawn for a particular bet.
 */
@Entity
@Table(name = "draws")
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class Draw {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @OneToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "bet_id", nullable = false, unique = true)
    private Bet bet;

    @Column(name = "die_one", nullable = false)
    private Integer dieOne;

    @Column(name = "die_two", nullable = false)
    private Integer dieTwo;

    @Column(name = "die_three", nullable = false)
    private Integer dieThree;

    /** Computed product of the three dice */
    @Column(name = "product_value", nullable = false)
    private Integer productValue;
}
