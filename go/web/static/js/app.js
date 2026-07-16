/**
 * ThreeDice Frontend Application
 *
 * Communicates with the Spring Boot REST API to let the player
 * place bets and view history.
 */
const API = '/api';
const PLAYER_ID = 1; // demo player seeded by DataInitializer

// ── DOM elements ──
const playerNameEl    = document.getElementById('player-name');
const playerBalanceEl = document.getElementById('player-balance');
const betForm         = document.getElementById('bet-form');
const betBtn          = document.getElementById('bet-btn');
const stakeInput      = document.getElementById('stake');
const predictedInput  = document.getElementById('predicted-value');

const resultSection   = document.getElementById('result-section');
const die1El          = document.getElementById('die1');
const die2El          = document.getElementById('die2');
const die3El          = document.getElementById('die3');
const productEl       = document.getElementById('product-value');
const statusEl        = document.getElementById('result-status');
const winningsEl      = document.getElementById('result-winnings');

const betHistoryBody  = document.getElementById('bet-history-body');
const txnHistoryBody  = document.getElementById('txn-history-body');

// ── Initialisation ──
document.addEventListener('DOMContentLoaded', () => {
    loadPlayer();
    loadBetHistory();
    loadTransactionHistory();
    initTabs();
    initQuickPicks();
    initQuickStakes();

    // Pre-highlight the default pick button
    const defaultPick = document.querySelector('.pick-btn[data-value="12"]');
    if (defaultPick) defaultPick.classList.add('active');
});

// ── Player Info ──
async function loadPlayer() {
    try {
        const res = await fetch(`${API}/players/${PLAYER_ID}`);
        if (!res.ok) throw new Error('Player not found');
        const data = await res.json();
        playerNameEl.textContent = data.username;
        updateBalance(data.balance);
    } catch (error) {
        console.error('Failed to load player:', error);
        showToast('Could not load player info');
    }
}

function updateBalance(amount) {
    playerBalanceEl.textContent = `$${Number(amount).toFixed(2)}`;
    playerBalanceEl.classList.add('updated');
    setTimeout(() => playerBalanceEl.classList.remove('updated'), 300);
}

// ── Place Bet ──
betForm.addEventListener('submit', async (e) => {
    e.preventDefault();

    const stake = Number.parseFloat(stakeInput.value);
    const predictedValue = Number.parseInt(predictedInput.value, 10);

    if (Number.isNaN(stake) || stake <= 0) return showToast('Enter a valid stake');
    if (Number.isNaN(predictedValue) || predictedValue < 1 || predictedValue > 216) {
        return showToast('Predicted value must be between 1 and 216');
    }

    betBtn.disabled = true;
    betBtn.textContent = 'Rolling...';

    try {
        const res = await fetch(`${API}/bets`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                playerId: PLAYER_ID,
                stake,
                predictedValue,
                idempotencyKey: crypto.randomUUID()
            })
        });

        if (!res.ok) {
            const err = await res.json();
            throw new Error(err.detail || 'Bet failed');
        }

        const data = await res.json();
        showResult(data);
        updateBalance(data.balanceAfter);
        loadBetHistory();
        loadTransactionHistory();

    } catch (err) {
        showToast(err.message);
    } finally {
        betBtn.disabled = false;
        betBtn.textContent = 'Roll the Dice!';
    }
});

function showResult(data) {
    resultSection.classList.remove('hidden');

    // Animate dice
    [die1El, die2El, die3El].forEach(el => {
        el.classList.remove('rolling');
        // Trigger reflow to restart animation
        el.offsetWidth; // eslint-disable-line no-unused-expressions
        el.classList.add('rolling');
    });

    setTimeout(() => {
        die1El.textContent = data.dieOne;
        die2El.textContent = data.dieTwo;
        die3El.textContent = data.dieThree;
        productEl.textContent = data.productValue;
    }, 200);

    const won = data.status === 'WON';
    statusEl.textContent = won ? '🎉 WON!' : '💔 LOST';
    statusEl.className = 'result-status ' + (won ? 'won' : 'lost');
    winningsEl.textContent = won ? `+$${Number(data.winnings).toFixed(2)}` : '';
}

// ── Bet History ──
async function loadBetHistory() {
    try {
        const res = await fetch(`${API}/bets/history/${PLAYER_ID}`);
        if (!res.ok) throw new Error('Failed to fetch bet history');
        const bets = await res.json();

        if (bets.length === 0) {
            betHistoryBody.innerHTML =
                '<tr><td colspan="7" class="empty-row">No bets yet. Place your first bet above!</td></tr>';
            return;
        }

        betHistoryBody.innerHTML = bets.map((b, i) => `
            <tr>
                <td>${bets.length - i}</td>
                <td>${b.predictedValue}</td>
                <td>${b.dieOne} × ${b.dieTwo} × ${b.dieThree}</td>
                <td>${b.productValue}</td>
                <td>$${Number(b.stake).toFixed(2)}</td>
                <td><span class="badge badge-${b.status.toLowerCase()}">${b.status}</span></td>
                <td>${formatTime(b.createdAt)}</td>
            </tr>
        `).join('');
    } catch {
        betHistoryBody.innerHTML =
            '<tr><td colspan="7" class="empty-row">Failed to load history</td></tr>';
    }
}

// ── Transaction History ──
async function loadTransactionHistory() {
    try {
        const res = await fetch(`${API}/transactions/history/${PLAYER_ID}`);
        if (!res.ok) throw new Error('Failed to fetch transaction history');
        const txns = await res.json();

        if (txns.length === 0) {
            txnHistoryBody.innerHTML =
                '<tr><td colspan="5" class="empty-row">No transactions yet.</td></tr>';
            return;
        }

        txnHistoryBody.innerHTML = txns.map((t, i) => `
            <tr>
                <td>${txns.length - i}</td>
                <td><span class="badge badge-${t.type.toLowerCase()}">${t.type}</span></td>
                <td>$${Number(t.amount).toFixed(2)}</td>
                <td>$${Number(t.balanceAfter).toFixed(2)}</td>
                <td>${formatTime(t.createdAt)}</td>
            </tr>
        `).join('');
    } catch {
        txnHistoryBody.innerHTML =
            '<tr><td colspan="5" class="empty-row">Failed to load transactions</td></tr>';
    }
}

// ── Tabs ──
function initTabs() {
    document.querySelectorAll('.tab').forEach(tab => {
        tab.addEventListener('click', () => {
            document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
            document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));
            tab.classList.add('active');
            document.getElementById(`tab-${tab.dataset.tab}`).classList.add('active');
        });
    });
}

// ── Quick Picks ──
function initQuickPicks() {
    document.querySelectorAll('.pick-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            // Set the predicted value input
            predictedInput.value = btn.dataset.value;
            // Toggle active state
            document.querySelectorAll('.pick-btn').forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            // Focus stake if empty
            if (!stakeInput.value) stakeInput.focus();
        });
    });

    // Clear active state when user types manually
    predictedInput.addEventListener('input', () => {
        document.querySelectorAll('.pick-btn').forEach(b => b.classList.remove('active'));
        const match = document.querySelector(`.pick-btn[data-value="${predictedInput.value}"]`);
        if (match) match.classList.add('active');
    });
}

// ── Quick Stakes ──
function initQuickStakes() {
    document.querySelectorAll('.stake-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            const amount = Number.parseFloat(btn.dataset.amount);
            const current = Number.parseFloat(stakeInput.value) || 0;
            stakeInput.value = (current + amount).toFixed(2);
        });
    });
}

// ── Utilities ──
function formatTime(iso) {
    if (!iso) return '-';
    const d = new Date(iso);
    return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })
         + ' ' + d.toLocaleDateString([], { month: 'short', day: 'numeric' });
}

function showToast(message) {
    // Remove existing toast if present
    const existing = document.querySelector('.toast');
    if (existing) existing.remove();

    const toast = document.createElement('div');
    toast.className = 'toast';
    toast.textContent = message;
    document.body.appendChild(toast);

    requestAnimationFrame(() => toast.classList.add('show'));
    setTimeout(() => {
        toast.classList.remove('show');
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}
