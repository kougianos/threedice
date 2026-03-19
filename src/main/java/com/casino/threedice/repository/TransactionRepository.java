package com.casino.threedice.repository;

import com.casino.threedice.entity.Transaction;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface TransactionRepository extends JpaRepository<Transaction, Long> {

    List<Transaction> findByClientIdOrderByCreatedAtDesc(Long clientId, Pageable pageable);
}
