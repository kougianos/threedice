package com.casino.threedice.repository;

import com.casino.threedice.entity.Client;
import jakarta.persistence.LockModeType;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Lock;
import org.springframework.data.jpa.repository.Query;

import java.util.Optional;

public interface ClientRepository extends JpaRepository<Client, Long> {

    /**
     * Loads the client with a pessimistic write lock (SELECT ... FOR UPDATE),
     * preventing concurrent transactions from reading stale balance data.
     */
    @Lock(LockModeType.PESSIMISTIC_WRITE)
    @Query("SELECT c FROM Client c WHERE c.id = :id")
    Optional<Client> findByIdForUpdate(Long id);
}
