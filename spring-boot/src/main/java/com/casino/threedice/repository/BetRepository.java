package com.casino.threedice.repository;

import com.casino.threedice.entity.Bet;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface BetRepository extends JpaRepository<Bet, Long> {

    List<Bet> findByClientIdOrderByCreatedAtDesc(Long clientId, Pageable pageable);

    boolean existsByIdempotencyKey(String idempotencyKey);
}
