package com.casino.threedice.exception;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.validation.ConstraintViolationException;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ProblemDetail;
import org.springframework.http.converter.HttpMessageNotReadableException;
import org.springframework.web.HttpRequestMethodNotSupportedException;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;
import org.springframework.web.method.annotation.MethodArgumentTypeMismatchException;
import org.springframework.web.servlet.resource.NoResourceFoundException;

import java.util.stream.Collectors;

@Slf4j
@RestControllerAdvice
public class GlobalExceptionHandler {

    @ExceptionHandler(PlayerNotFoundException.class)
    ProblemDetail handlePlayerNotFound(PlayerNotFoundException ex) {
        log.warn(ex.getMessage());
        return ProblemDetail.forStatusAndDetail(HttpStatus.NOT_FOUND, ex.getMessage());
    }

    @ExceptionHandler(InsufficientBalanceException.class)
    ProblemDetail handleInsufficientBalance(InsufficientBalanceException ex) {
        log.warn(ex.getMessage());
        return ProblemDetail.forStatusAndDetail(HttpStatus.BAD_REQUEST, ex.getMessage());
    }

    @ExceptionHandler(DuplicateBetException.class)
    ProblemDetail handleDuplicateBet(DuplicateBetException ex) {
        log.warn(ex.getMessage());
        return ProblemDetail.forStatusAndDetail(HttpStatus.CONFLICT, ex.getMessage());
    }

    @ExceptionHandler(MethodArgumentNotValidException.class)
    ProblemDetail handleValidation(MethodArgumentNotValidException ex) {
        String errors = ex.getBindingResult().getFieldErrors().stream()
                .map(fe -> fe.getField() + ": " + fe.getDefaultMessage())
                .sorted()
                .collect(Collectors.joining("; "));
        log.warn("Validation failed: {}", errors);
        return ProblemDetail.forStatusAndDetail(HttpStatus.BAD_REQUEST, errors);
    }

    @ExceptionHandler(ConstraintViolationException.class)
    ProblemDetail handleConstraintViolation(ConstraintViolationException ex) {
        String errors = ex.getConstraintViolations().stream()
                .map(v -> v.getPropertyPath() + ": " + v.getMessage())
                .sorted()
                .collect(Collectors.joining("; "));
        log.warn("Constraint violation: {}", errors);
        return ProblemDetail.forStatusAndDetail(HttpStatus.BAD_REQUEST, errors);
    }

    /**
     * A path variable that will not convert, e.g. /api/players/abc.
     * Without this the catch-all below would report a client mistake as a 500.
     */
    @ExceptionHandler(MethodArgumentTypeMismatchException.class)
    ProblemDetail handleTypeMismatch(MethodArgumentTypeMismatchException ex) {
        String detail = "%s must be a number: %s".formatted(ex.getName(), ex.getValue());
        log.warn(detail);
        return ProblemDetail.forStatusAndDetail(HttpStatus.BAD_REQUEST, detail);
    }

    /**
     * A body that is absent or not parseable as JSON.
     */
    @ExceptionHandler(HttpMessageNotReadableException.class)
    ProblemDetail handleUnreadableBody(HttpMessageNotReadableException ex) {
        log.warn("Malformed request body: {}", ex.getMessage());
        return ProblemDetail.forStatusAndDetail(HttpStatus.BAD_REQUEST, "Malformed JSON request body");
    }

    /**
     * An unmapped path, which would otherwise fall through to the catch-all as a 500.
     */
    @ExceptionHandler(NoResourceFoundException.class)
    ProblemDetail handleNoResource(NoResourceFoundException ex) {
        String detail = "No endpoint %s /%s".formatted(ex.getHttpMethod(), ex.getResourcePath());
        log.warn(detail);
        return ProblemDetail.forStatusAndDetail(HttpStatus.NOT_FOUND, detail);
    }

    /**
     * A known path called with the wrong verb, e.g. POST /api/players/1.
     */
    @ExceptionHandler(HttpRequestMethodNotSupportedException.class)
    ProblemDetail handleMethodNotSupported(HttpRequestMethodNotSupportedException ex,
                                           HttpServletRequest request) {
        String detail = "Method %s is not supported for %s"
                .formatted(ex.getMethod(), request.getRequestURI());
        log.warn(detail);
        return ProblemDetail.forStatusAndDetail(HttpStatus.METHOD_NOT_ALLOWED, detail);
    }

    @ExceptionHandler(Exception.class)
    ProblemDetail handleGeneric(Exception ex) {
        log.error("Unexpected error", ex);
        return ProblemDetail.forStatusAndDetail(
                HttpStatus.INTERNAL_SERVER_ERROR, "An unexpected error occurred");
    }
}
