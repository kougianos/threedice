/**
 * ThreeDice Frontend Application
 *
 * Talks to the ThreeDice REST API, which the Spring Boot and Go services
 * implement identically. These files ship byte for byte in both, so nothing
 * here may assume one runtime over the other.
 */
const API = '/api';
const PLAYER_ID = 1; // the demo player each service seeds on first startup

// ── Game facts ──
// Three six-sided dice produce 216 ordered rolls but only 40 distinct products.
// Counting them here, rather than transcribing percentages by hand, is what
// keeps the odds on screen in step with the ones the service actually pays.
const TOTAL_ROLLS = 216;

const COMBINATIONS = (() => {
    const counts = new Map();
    for (let a = 1; a <= 6; a++) {
        for (let b = 1; b <= 6; b++) {
            for (let c = 1; c <= 6; c++) {
                const product = a * b * c;
                counts.set(product, (counts.get(product) || 0) + 1);
            }
        }
    }
    return counts;
})();

const VALID_PRODUCTS = [...COMBINATIONS.keys()].sort((a, b) => a - b);

const MIN_STAKE = 1;
const MAX_STAKE = 10000;

/** The payout bands the services apply, mirrored so the UI can quote them. */
function oddsFor(product) {
    if (product < 9) return 2;
    if (product < 120) return 5;
    return 2;
}

function chanceOf(product) {
    return ((COMBINATIONS.get(product) || 0) / TOTAL_ROLLS) * 100;
}

function isValidProduct(value) {
    return COMBINATIONS.has(value);
}

/** The reachable products either side of an unreachable guess. */
function nearestValidProducts(value) {
    const below = [...VALID_PRODUCTS].reverse().find((p) => p < value);
    const above = VALID_PRODUCTS.find((p) => p > value);
    return [below, above].filter((p) => p !== undefined);
}

// The eight products offered as one-tap picks.
const QUICK_PICKS = [6, 12, 24, 36, 8, 30, 60, 120];

// ── DOM elements ──
const playerNameEl    = document.getElementById('player-name');
const playerBalanceEl = document.getElementById('player-balance');
const betForm         = document.getElementById('bet-form');
const betBtn          = document.getElementById('bet-btn');
const stakeInput      = document.getElementById('stake');
const stakeError      = document.getElementById('stake-error');
const predictedInput  = document.getElementById('predicted-value');
const predictedError  = document.getElementById('predicted-error');
const quickPicksEl    = document.getElementById('quick-picks');

const diceDisplay     = document.getElementById('dice-display');
const die1El          = document.getElementById('die1');
const die2El          = document.getElementById('die2');
const die3El          = document.getElementById('die3');
const productEl       = document.getElementById('product-value');
const predictedEchoEl = document.getElementById('predicted-echo');
const resultDetails   = document.getElementById('result-details');
const statusEl        = document.getElementById('result-status');
const winningsEl      = document.getElementById('result-winnings');
const restingEl       = document.getElementById('result-resting');

const betHistoryBody  = document.getElementById('bet-history-body');
const txnHistoryBody  = document.getElementById('txn-history-body');
const announcer       = document.getElementById('announcer');

const dice = [die1El, die2El, die3El];

let currentBalance = null;

// ── Initialisation ──
document.addEventListener('DOMContentLoaded', () => {
    renderQuickPicks();
    loadPlayer();
    loadBetHistory();
    loadTransactionHistory();
    initTabs();
    initQuickStakes();
    syncPredictionEcho();
});

// ── Announcements ──
/** Assistive tech has no other way to learn a roll happened. */
function announce(message) {
    announcer.textContent = message;
}

// ── Player Info ──
async function loadPlayer() {
    try {
        const res = await fetch(`${API}/players/${PLAYER_ID}`);
        if (!res.ok) throw new Error('Player not found');
        const data = await res.json();
        playerNameEl.textContent = data.username;
        currentBalance = Number(data.balance);
        playerBalanceEl.textContent = formatMoney(currentBalance);
    } catch {
        showToast('Could not load your player details. Refresh to try again.');
    }
}

function updateBalance(amount) {
    const next = Number(amount);
    const direction = currentBalance === null || next === currentBalance
        ? null
        : (next > currentBalance ? 'credited' : 'debited');

    currentBalance = next;
    playerBalanceEl.textContent = formatMoney(next);

    // The pill used to play the same celebratory bounce whichever way the
    // money went. It now marks the direction it actually moved.
    playerBalanceEl.classList.remove('credited', 'debited');
    if (direction) {
        void playerBalanceEl.offsetWidth;
        playerBalanceEl.classList.add(direction);
        setTimeout(() => playerBalanceEl.classList.remove(direction), 400);
    }
}

// ── Quick Picks ──
function renderQuickPicks() {
    const expectedValue = (p) => (chanceOf(p) / 100) * oddsFor(p);
    const bestReturn = Math.max(...QUICK_PICKS.map(expectedValue));

    quickPicksEl.innerHTML = QUICK_PICKS.map((value) => {
        const chance = chanceOf(value).toFixed(1);
        const odds = oddsFor(value);
        // 12 and 24 tie exactly, so "best" is not a single button.
        const isBest = Math.abs(expectedValue(value) - bestReturn) < 1e-9;
        const label = `Predict ${value}. ${chance} percent chance, pays ${odds}x`
            + (isBest ? ', the best return on offer here' : '');

        return `<button type="button" class="pick-btn${isBest ? ' best' : ''}"
                        data-value="${value}" aria-pressed="false" aria-label="${label}">
                    <span class="pick-value">${value}</span>
                    <span class="pick-meta">${chance}% &middot; ${odds}x</span>
                </button>`;
    }).join('');

    quickPicksEl.querySelectorAll('.pick-btn').forEach((btn) => {
        btn.addEventListener('click', () => {
            predictedInput.value = btn.dataset.value;
            clearFieldError(predictedInput, predictedError);
            syncPredictionEcho();
        });
    });

    syncPickSelection();
}

function syncPickSelection() {
    const current = predictedInput.value.trim();
    quickPicksEl.querySelectorAll('.pick-btn').forEach((btn) => {
        btn.setAttribute('aria-pressed', String(btn.dataset.value === current));
    });
}

// ── Quick Stakes ──
function initQuickStakes() {
    document.querySelectorAll('.stake-btn').forEach((btn) => {
        btn.addEventListener('click', () => {
            const amount = Number.parseFloat(btn.dataset.amount);
            const current = Number.parseFloat(stakeInput.value) || 0;
            stakeInput.value = Math.min(current + amount, MAX_STAKE).toFixed(2);
            clearFieldError(stakeInput, stakeError);
        });
    });
}

// ── Prediction echo ──
/** Keeps the other half of the comparison visible, before and after a roll. */
function syncPredictionEcho() {
    const raw = predictedInput.value.trim();
    predictedEchoEl.textContent = raw === '' ? 'no prediction yet' : `you predicted ${raw}`;
    syncPickSelection();
}

predictedInput.addEventListener('input', () => {
    clearFieldError(predictedInput, predictedError);
    syncPredictionEcho();
});

stakeInput.addEventListener('input', () => clearFieldError(stakeInput, stakeError));

// ── Field validation ──
function setFieldError(input, target, message) {
    input.setAttribute('aria-invalid', 'true');
    target.textContent = message;
    target.hidden = false;
}

function clearFieldError(input, target) {
    input.removeAttribute('aria-invalid');
    target.textContent = '';
    target.hidden = true;
}

/**
 * Answers at the field, before anything is sent. The services reject the same
 * cases, but a round trip is a poor way to learn that 13 is not reachable.
 */
function validateBet() {
    clearFieldError(stakeInput, stakeError);
    clearFieldError(predictedInput, predictedError);

    const errors = [];
    const stake = Number.parseFloat(stakeInput.value);
    const predicted = Number.parseInt(predictedInput.value, 10);

    if (!Number.isFinite(stake)) {
        errors.push([stakeInput, stakeError, 'Enter a stake as a number, like 10.00.']);
    } else if (stake < MIN_STAKE) {
        errors.push([stakeInput, stakeError, `The smallest stake is ${formatMoney(MIN_STAKE)}.`]);
    } else if (stake > MAX_STAKE) {
        errors.push([stakeInput, stakeError, `The largest stake is ${formatMoney(MAX_STAKE)}.`]);
    } else if (currentBalance !== null && stake > currentBalance) {
        errors.push([stakeInput, stakeError,
            `That is more than your balance of ${formatMoney(currentBalance)}.`]);
    }

    if (!Number.isInteger(predicted)) {
        errors.push([predictedInput, predictedError, 'Enter a whole number.']);
    } else if (!isValidProduct(predicted)) {
        const nearby = nearestValidProducts(predicted);
        const suggestion = nearby.length ? ` Try ${nearby.join(' or ')}.` : '';
        errors.push([predictedInput, predictedError,
            `Three dice cannot make ${predicted}.${suggestion}`]);
    }

    errors.forEach(([input, target, message]) => setFieldError(input, target, message));

    if (errors.length) {
        const [firstInput, , firstMessage] = errors[0];
        firstInput.focus();
        announce(firstMessage);
        return null;
    }

    return { stake, predicted };
}

// ── Place Bet ──
const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)');

betForm.addEventListener('submit', async (e) => {
    e.preventDefault();

    const bet = validateBet();
    if (!bet) return;

    betBtn.disabled = true;
    betBtn.textContent = 'Rolling...';
    startRoll();

    try {
        const res = await fetch(`${API}/bets`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                playerId: PLAYER_ID,
                stake: bet.stake,
                predictedValue: bet.predicted,
                idempotencyKey: crypto.randomUUID()
            })
        });

        if (!res.ok) {
            const problem = await res.json().catch(() => ({}));
            throw new Error(readableDetail(problem.detail) || 'The bet could not be placed.');
        }

        const data = await res.json();
        // The balance and the history are part of the outcome too, so they wait
        // for the dice along with everything else.
        await revealResult(data, bet.predicted);
        loadBetHistory();
        loadTransactionHistory();
    } catch (err) {
        cancelRoll();
        showToast(err.message);
        announce(err.message);
    } finally {
        betBtn.disabled = false;
        betBtn.textContent = 'Roll the Dice';
    }
});

/** Clears the previous bet so it cannot masquerade as this one. */
function startRoll() {
    dice.forEach((el) => {
        el.classList.remove('rolling');
        el.textContent = '?';
    });
    productEl.textContent = '?';
    resultDetails.hidden = true;
    statusEl.className = 'result-status';
    statusEl.textContent = '';
    winningsEl.textContent = '';
    predictedEchoEl.classList.remove('matched');
    restingEl.hidden = true;
    diceDisplay.classList.add('pending');
}

function cancelRoll() {
    diceDisplay.classList.remove('pending');
    restingEl.hidden = false;
    restingEl.textContent = 'That roll did not go through. Your stake was not taken.';
}

/**
 * Resolves once the outcome is on screen. Nothing about the result appears
 * before the dice that produce it: the verdict used to be written 200ms early,
 * and the balance used to move the moment the response landed.
 */
function revealResult(data, predicted) {
    diceDisplay.classList.remove('pending');
    restingEl.hidden = true;

    dice.forEach((el) => {
        el.classList.remove('rolling');
        void el.offsetWidth; // restart the animation
        el.classList.add('rolling');
    });

    const won = data.status === 'WON';

    return new Promise((resolve) => {
        const settle = () => {
            die1El.textContent = data.dieOne;
            die2El.textContent = data.dieTwo;
            die3El.textContent = data.dieThree;
            productEl.textContent = data.productValue;

            predictedEchoEl.textContent = `you predicted ${predicted}`;
            predictedEchoEl.classList.toggle('matched', won);

            statusEl.textContent = won ? 'WON' : 'LOST';
            statusEl.className = 'result-status ' + (won ? 'won' : 'lost');
            winningsEl.textContent = won ? `+${formatMoney(data.winnings)}` : '';
            resultDetails.hidden = false;

            updateBalance(data.balanceAfter);

            announce(
                `Rolled ${data.dieOne}, ${data.dieTwo} and ${data.dieThree}. `
                + `Product ${data.productValue}, you predicted ${predicted}. `
                + (won ? `Won ${formatMoney(data.winnings)}. ` : 'Lost. ')
                + `Balance ${formatMoney(data.balanceAfter)}.`
            );
            resolve();
        };

        if (prefersReducedMotion.matches) settle();
        else setTimeout(settle, 200);
    });
}

// ── Bet History ──
async function loadBetHistory() {
    try {
        const res = await fetch(`${API}/bets/history/${PLAYER_ID}`);
        if (!res.ok) throw new Error('Failed to fetch bet history');
        const bets = await res.json();

        if (bets.length === 0) {
            betHistoryBody.innerHTML =
                '<tr><td colspan="7" class="empty-row">No bets yet. Place your first bet above.</td></tr>';
            return;
        }

        betHistoryBody.innerHTML = bets.map((b) => `
            <tr>
                <td class="num" data-label="Bet">${Number(b.betId)}</td>
                <td class="num" data-label="Predicted">${Number(b.predictedValue)}</td>
                <td class="num" data-label="Product">${Number(b.productValue)}</td>
                <td data-label="Dice">${Number(b.dieOne)} &times; ${Number(b.dieTwo)} &times; ${Number(b.dieThree)}</td>
                <td class="num" data-label="Stake">${formatMoney(b.stake)}</td>
                <td data-label="Result">${badge(b.status)}</td>
                <td class="col-time" data-label="Time" title="${absoluteTime(b.createdAt)}">${relativeTime(b.createdAt)}</td>
            </tr>
        `).join('');
    } catch {
        betHistoryBody.innerHTML =
            '<tr><td colspan="7" class="empty-row">Could not load your bets.</td></tr>';
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

        txnHistoryBody.innerHTML = txns.map((t) => `
            <tr>
                <td class="num" data-label="Txn">${Number(t.transactionId)}</td>
                <td data-label="Type">${badge(t.type)}</td>
                <td class="num" data-label="Amount">${formatMoney(t.amount)}</td>
                <td class="num" data-label="Balance after">${formatMoney(t.balanceAfter)}</td>
                <td class="col-time" data-label="Time" title="${absoluteTime(t.createdAt)}">${relativeTime(t.createdAt)}</td>
            </tr>
        `).join('');
    } catch {
        txnHistoryBody.innerHTML =
            '<tr><td colspan="5" class="empty-row">Could not load your transactions.</td></tr>';
    }
}

function badge(value) {
    const key = String(value).toLowerCase().replace(/[^a-z]/g, '');
    return `<span class="badge badge-${key}">${key.toUpperCase()}</span>`;
}

// ── Tabs ──
function initTabs() {
    const tabs = [...document.querySelectorAll('[role="tab"]')];

    const select = (tab) => {
        tabs.forEach((t) => {
            const selected = t === tab;
            t.setAttribute('aria-selected', String(selected));
            t.tabIndex = selected ? 0 : -1;
            t.classList.toggle('active', selected);
            const panel = document.getElementById(t.getAttribute('aria-controls'));
            panel.hidden = !selected;
            panel.classList.toggle('active', selected);
        });
    };

    tabs.forEach((tab, i) => {
        tab.addEventListener('click', () => select(tab));
        tab.addEventListener('keydown', (e) => {
            if (e.key !== 'ArrowRight' && e.key !== 'ArrowLeft') return;
            e.preventDefault();
            const next = tabs[(i + (e.key === 'ArrowRight' ? 1 : tabs.length - 1)) % tabs.length];
            select(next);
            next.focus();
        });
    });
}

// ── Utilities ──
const MONEY = new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' });

/** Always two decimals and a grouped thousand, whatever scale the service sent. */
function formatMoney(amount) {
    return MONEY.format(Number(amount));
}

/** Strips the `field: ` prefix the services put on validation details. */
function readableDetail(detail) {
    if (!detail) return '';
    return String(detail)
        .split(', ')
        .map((part) => part.replace(/^[a-zA-Z]+:\s*/, ''))
        .join('. ');
}

function absoluteTime(iso) {
    if (!iso) return '';
    return new Date(iso).toLocaleString();
}

function relativeTime(iso) {
    if (!iso) return '-';
    const then = new Date(iso).getTime();
    const seconds = Math.round((Date.now() - then) / 1000);

    if (seconds < 10) return 'just now';
    if (seconds < 60) return `${seconds}s ago`;
    if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
    if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
    return new Date(iso).toLocaleDateString([], { month: 'short', day: 'numeric' });
}

/**
 * Server and network failures only; validation answers at the field. The show
 * class is applied after a forced reflow rather than inside requestAnimationFrame,
 * which never fires while the tab is in the background.
 */
function showToast(message) {
    document.querySelector('.toast')?.remove();

    const toast = document.createElement('div');
    toast.className = 'toast';
    toast.setAttribute('role', 'alert');

    const text = document.createElement('span');
    text.textContent = message;

    const dismiss = document.createElement('button');
    dismiss.type = 'button';
    dismiss.className = 'toast-dismiss';
    dismiss.setAttribute('aria-label', 'Dismiss this message');
    dismiss.textContent = '×';

    let timer;
    const close = () => {
        clearTimeout(timer);
        toast.classList.remove('show');
        setTimeout(() => toast.remove(), 300);
    };

    dismiss.addEventListener('click', close);
    toast.append(text, dismiss);
    document.body.appendChild(toast);

    void toast.offsetWidth;
    toast.classList.add('show');
    timer = setTimeout(close, 8000);
}
