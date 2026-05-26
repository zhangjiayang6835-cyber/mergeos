<template>
  <div v-if="!user" class="auth-shell">
    <section class="auth-panel">
      <div class="brand-lockup">
        <div class="brand-mark">
          <PanelsTopLeft :size="25" />
        </div>
        <div>
          <p class="eyebrow">Client portal</p>
          <h1>MergeOS</h1>
        </div>
      </div>

      <div class="auth-copy">
        <h2>Fund a website, get a private bounty repo, track every paid task.</h2>
        <p>Register as a client, verify PayPal or crypto payment, and MergeOS converts the budget into MERGE credits for human and agent delivery.</p>
      </div>

      <div class="segmented">
        <button :class="{ active: authMode === 'register' }" @click="authMode = 'register'">Register</button>
        <button :class="{ active: authMode === 'login' }" @click="authMode = 'login'">Login</button>
      </div>

      <form class="auth-form" @submit.prevent="submitAuth">
        <label v-if="authMode === 'register'">
          Full name
          <input v-model="authForm.name" autocomplete="name" />
        </label>
        <label v-if="authMode === 'register'">
          Company
          <input v-model="authForm.company_name" autocomplete="organization" />
        </label>
        <label>
          Email
          <input v-model="authForm.email" autocomplete="email" type="email" />
        </label>
        <label>
          Password
          <input v-model="authForm.password" autocomplete="current-password" type="password" />
        </label>
        <button class="primary-button" :disabled="authBusy">
          <LogIn :size="17" />
          <span>{{ authBusy ? 'Working...' : authMode === 'register' ? 'Create account' : 'Login' }}</span>
        </button>
        <p v-if="errorMessage" class="error-line">{{ errorMessage }}</p>
      </form>

      <div class="auth-runtime">
        <span>{{ runtimeConfig?.payment_mode || 'loading payment' }}</span>
        <span>{{ runtimeConfig?.repo_provider || 'loading repo' }}</span>
        <span>{{ runtimeConfig?.smtp_ready ? 'smtp ready' : 'email log mode' }}</span>
      </div>
    </section>
  </div>

  <div v-else class="app-shell">
    <header class="topbar">
      <div class="brand-lockup compact">
        <div class="brand-mark">
          <PanelsTopLeft :size="23" />
        </div>
        <div>
          <p class="eyebrow">Private client workspace</p>
          <h1>MergeOS</h1>
        </div>
      </div>
      <nav class="top-tabs">
        <button :class="{ active: portalTab === 'workspace' }" @click="portalTab = 'workspace'">
          <LayoutDashboard :size="16" />
          <span>Workspace</span>
        </button>
        <button :class="{ active: portalTab === 'billing' }" @click="portalTab = 'billing'">
          <WalletCards :size="16" />
          <span>Billing</span>
        </button>
        <button :class="{ active: portalTab === 'inbox' }" @click="portalTab = 'inbox'">
          <Mail :size="16" />
          <span>Email</span>
        </button>
      </nav>
      <div class="topbar-spacer" />
      <div class="status-pill">
        <ShieldCheck :size="16" />
        <span>{{ statusLabel }}</span>
      </div>
      <button class="icon-button" title="Refresh workspace" @click="refreshAll">
        <RefreshCw :size="18" />
      </button>
      <button class="icon-button" title="Logout" @click="logout">
        <LogOut :size="18" />
      </button>
    </header>

    <aside class="customer-panel">
      <div class="panel-heading">
        <UserRound :size="18" />
        <span>Customer</span>
      </div>

      <div class="profile-card">
        <strong>{{ user.name }}</strong>
        <span>{{ user.company_name || 'Independent client' }}</span>
        <small>{{ user.email }}</small>
      </div>

      <div class="runtime-card">
        <div>
          <span>Payment</span>
          <strong>{{ runtimeConfig?.payment_mode || 'loading' }}</strong>
        </div>
        <div>
          <span>Repo</span>
          <strong>{{ runtimeConfig?.repo_provider || 'loading' }}</strong>
        </div>
        <div>
          <span>Email</span>
          <strong>{{ runtimeConfig?.smtp_ready ? 'smtp' : 'log' }}</strong>
        </div>
        <div>
          <span>Token</span>
          <strong>{{ tokenSymbol }}</strong>
        </div>
      </div>

      <label>
        Contact name
        <input v-model="projectForm.client_name" />
      </label>
      <label>
        Contact email
        <input v-model="projectForm.client_email" type="email" />
      </label>
      <label>
        Phone
        <input v-model="projectForm.phone" />
      </label>
      <label>
        Company
        <input v-model="projectForm.company_name" />
      </label>
    </aside>

    <main class="portal-main">
      <section class="summary-strip">
        <article>
          <span>Total funded</span>
          <strong>{{ money(totalBudget) }}</strong>
        </article>
        <article>
          <span>MERGE reserved</span>
          <strong>{{ money(totalPool) }}</strong>
        </article>
        <article>
          <span>Open tasks</span>
          <strong>{{ openTasks.length }}</strong>
        </article>
        <article>
          <span>Paid tasks</span>
          <strong>{{ acceptedTasks.length }}</strong>
        </article>
      </section>

      <section v-if="portalTab === 'workspace'" class="workspace-grid">
        <div class="canvas-column">
          <div class="canvas-toolbar">
            <div>
              <p class="eyebrow">Website order</p>
              <h2>{{ currentProject?.title || 'Create a funded website project' }}</h2>
            </div>
            <div class="toolbar-metrics">
              <span>{{ currentTasks.length }} issues</span>
              <span>{{ currentProject?.repo_provider || runtimeConfig?.repo_provider || 'local-git' }}</span>
            </div>
          </div>

          <section class="builder-canvas">
            <div class="canvas-section hero-section">
              <div class="section-handle">BRIEF</div>
              <div>
                <p class="eyebrow">{{ projectForm.site_type }} / {{ projectForm.timeline }}</p>
                <h3>{{ currentProject?.company_name || projectForm.company_name }} delivery room</h3>
                <p>{{ currentProject?.brief || projectForm.brief }}</p>
              </div>
              <div class="quote-block">
                <span>{{ currentProject?.payment_provider || 'checkout pending' }}</span>
                <strong>{{ money(currentProject?.budget_cents || projectForm.budget_cents) }}</strong>
              </div>
            </div>

            <div v-if="currentTasks.length" class="canvas-grid">
              <button
                v-for="task in currentTasks"
                :key="task.id"
                :class="['task-tile', { selected: selectedTask?.id === task.id, accepted: task.status === 'accepted' }]"
                @click="selectTask(task)"
              >
                <span class="issue-label">Issue #{{ task.issue_number }}</span>
                <strong>{{ task.title }}</strong>
                <small>{{ task.required_worker_kind }} / {{ money(task.reward_cents) }} {{ tokenSymbol }}</small>
              </button>
            </div>

            <div v-else class="empty-canvas">
              <Database :size="30" />
              <strong>No funded bounty yet</strong>
              <span>{{ runtimeConfig?.dev_payment_enabled ? `Use ${runtimeConfig.dev_payment_code} as the local payment reference.` : 'Configure PayPal or crypto credentials first.' }}</span>
            </div>
          </section>

          <section class="repo-strip">
            <div class="strip-header">
              <GitBranch :size="18" />
              <a v-if="currentProject?.repo_url" :href="currentProject.repo_url" target="_blank" rel="noreferrer">
                {{ currentProject.bounty_repo_name }}
              </a>
              <span v-else>mergeos-bounties repo pending</span>
            </div>
            <div class="task-list">
              <button
                v-for="task in currentTasks"
                :key="task.id"
                :class="['task-row', { selected: selectedTask?.id === task.id }]"
                @click="selectTask(task)"
              >
                <span :class="['status-dot', task.status]" />
                <span>{{ task.title }}</span>
                <strong>{{ money(task.reward_cents) }}</strong>
              </button>
            </div>
          </section>
        </div>

        <aside class="order-panel">
          <div class="panel-heading">
            <FilePlus2 :size="18" />
            <span>New project</span>
          </div>
          <label>
            Project name
            <input v-model="projectForm.title" />
          </label>
          <label>
            Site type
            <select v-model="projectForm.site_type">
              <option>Landing page</option>
              <option>Business website</option>
              <option>SaaS website</option>
              <option>Storefront</option>
              <option>Web app shell</option>
            </select>
          </label>
          <label>
            Package
            <select v-model="projectForm.package_tier">
              <option>Launch</option>
              <option>Growth</option>
              <option>Scale</option>
            </select>
          </label>
          <label>
            Timeline
            <select v-model="projectForm.timeline">
              <option>7 days</option>
              <option>14 days</option>
              <option>30 days</option>
            </select>
          </label>
          <label>
            Budget USD
            <input v-model.number="budgetUsd" min="100" type="number" />
          </label>
          <label>
            Brief
            <textarea v-model="projectForm.brief" rows="5" />
          </label>
          <button class="primary-button" :disabled="creating" @click="createProject">
            <CreditCard :size="17" />
            <span>{{ creating ? 'Funding...' : 'Verify payment and create repo' }}</span>
          </button>
          <p v-if="errorMessage" class="error-line">{{ errorMessage }}</p>
        </aside>
      </section>

      <section v-if="portalTab === 'billing'" class="billing-grid">
        <div class="checkout-panel">
          <div class="panel-heading">
            <WalletCards :size="18" />
            <span>Checkout</span>
          </div>
          <div class="payment-choice">
            <button :class="{ active: projectForm.payment_method === 'paypal' }" @click="projectForm.payment_method = 'paypal'">PayPal</button>
            <button :class="{ active: projectForm.payment_method === 'crypto' }" @click="projectForm.payment_method = 'crypto'">Crypto</button>
          </div>
          <button
            v-if="projectForm.payment_method === 'paypal'"
            class="secondary-button"
            :disabled="!runtimeConfig?.paypal_ready || preparingPayPal"
            @click="preparePayPalOrder"
          >
            <WalletCards :size="17" />
            <span>{{ preparingPayPal ? 'Creating order...' : 'Create PayPal order' }}</span>
          </button>
          <a v-if="paypalOrder?.approval_url" class="approval-link" :href="paypalOrder.approval_url" target="_blank" rel="noreferrer">
            <ExternalLink :size="16" />
            <span>Open PayPal approval</span>
          </a>
          <label>
            Payment reference
            <input v-model="projectForm.payment_reference" :placeholder="paymentReferencePlaceholder" />
          </label>
          <div v-if="projectForm.payment_method === 'crypto'" class="receiver-card">
            <span>Receiver</span>
            <strong>{{ runtimeConfig?.crypto_receiver || 'Configure CRYPTO_RECEIVER in backend env' }}</strong>
          </div>
          <div class="billing-ledger">
            <div v-for="entry in recentLedger" :key="entry.sequence" class="ledger-line">
              <span>#{{ entry.sequence }} {{ entry.type }}</span>
              <strong>{{ money(entry.amount_cents) }}</strong>
            </div>
          </div>
        </div>

        <div class="project-list">
          <div class="panel-heading">
            <FolderKanban :size="18" />
            <span>Funded projects</span>
          </div>
          <button
            v-for="project in projects"
            :key="project.id"
            :class="['project-row', { selected: currentProject?.id === project.id }]"
            @click="selectProject(project)"
          >
            <span>
              <strong>{{ project.title }}</strong>
              <small>{{ project.bounty_repo_name }}</small>
            </span>
            <b>{{ money(project.budget_cents) }}</b>
          </button>
        </div>
      </section>

      <section v-if="portalTab === 'inbox'" class="inbox-grid">
        <div class="email-list">
          <div class="panel-heading">
            <Mail :size="18" />
            <span>Customer emails</span>
          </div>
          <article v-for="note in notifications" :key="note.id" class="email-card">
            <span>{{ note.status }}</span>
            <strong>{{ note.subject }}</strong>
            <p>{{ note.body }}</p>
          </article>
        </div>
      </section>
    </main>

    <aside class="inspector">
      <div class="panel-heading">
        <SplitSquareVertical :size="18" />
        <span>Task inspector</span>
      </div>

      <div class="repo-summary">
        <p class="eyebrow">Child bounty repo</p>
        <h3>{{ currentProject?.bounty_repo_name || 'Not created' }}</h3>
        <p>{{ currentProject?.repo_provider || runtimeConfig?.repo_provider || 'local-git' }}</p>
        <a v-if="currentProject?.repo_url" :href="currentProject.repo_url" target="_blank" rel="noreferrer">
          <ExternalLink :size="16" />
          <span>Open repo</span>
        </a>
      </div>

      <div v-if="selectedTask" class="task-inspector">
        <p class="eyebrow">Selected issue</p>
        <h3>{{ selectedTask.title }}</h3>
        <p>{{ selectedTask.acceptance }}</p>
        <a v-if="selectedTask.issue_url" :href="selectedTask.issue_url" target="_blank" rel="noreferrer">
          <ExternalLink :size="16" />
          <span>Open issue</span>
        </a>

        <div class="manifest-grid">
          <span>Required</span>
          <strong>{{ selectedTask.required_worker_kind }}</strong>
          <span>Suggested</span>
          <strong>{{ selectedTask.suggested_agent_type || 'human-review' }}</strong>
          <span>Reward</span>
          <strong>{{ money(selectedTask.reward_cents) }} {{ tokenSymbol }}</strong>
          <span>Proof</span>
          <strong>{{ selectedTask.proof_hash ? shortHash(selectedTask.proof_hash) : 'pending' }}</strong>
        </div>

        <label>
          Worker kind
          <select v-model="workerForm.worker_kind">
            <option value="human">Human</option>
            <option value="agent">Agent</option>
            <option value="hybrid">Hybrid</option>
          </select>
        </label>
        <label>
          Worker ID
          <input v-model="workerForm.worker_id" placeholder="github:alice or agent:web-001" />
        </label>
        <label>
          Agent type
          <input v-model="workerForm.agent_type" :disabled="workerForm.worker_kind === 'human'" placeholder="frontend-agent" />
        </label>
        <button class="primary-button" :disabled="selectedTask.status === 'accepted' || accepting" @click="acceptSelectedTask">
          <CheckCircle2 :size="17" />
          <span>{{ selectedTask.status === 'accepted' ? 'Paid' : 'Accept and pay' }}</span>
        </button>
      </div>
    </aside>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue';
import {
  CheckCircle2,
  CreditCard,
  Database,
  ExternalLink,
  FilePlus2,
  FolderKanban,
  GitBranch,
  LayoutDashboard,
  LogIn,
  LogOut,
  Mail,
  PanelsTopLeft,
  RefreshCw,
  ShieldCheck,
  SplitSquareVertical,
  UserRound,
  WalletCards,
} from '@lucide/vue/dist/esm/lucide-vue.mjs';

const runtimeConfig = ref(null);
const user = ref(null);
const authMode = ref('register');
const token = ref(localStorage.getItem('mergeos_token') || '');
const authBusy = ref(false);
const projects = ref([]);
const tasks = ref([]);
const ledger = ref([]);
const notifications = ref([]);
const selectedProjectId = ref('');
const selectedTask = ref(null);
const portalTab = ref('workspace');
const creating = ref(false);
const accepting = ref(false);
const preparingPayPal = ref(false);
const errorMessage = ref('');
const paypalOrder = ref(null);

const authForm = reactive({
  name: 'Thanh Truc Client',
  company_name: 'Thanh Truc Solutions',
  email: 'client@mergeos.local',
  password: 'mergeos123',
});

const projectForm = reactive({
  client_name: 'Thanh Truc Client',
  company_name: 'Thanh Truc Solutions',
  client_email: 'client@mergeos.local',
  phone: '+84 900 000 000',
  title: 'Elementor-style company website',
  site_type: 'Business website',
  package_tier: 'Growth',
  timeline: '14 days',
  brief: 'Build a polished responsive website with services, portfolio, lead form, payment-ready checkout, and customer dashboard preview.',
  budget_cents: 240000,
  payment_method: 'paypal',
  payment_reference: 'LOCAL-PAID',
});

const workerForm = reactive({
  worker_kind: 'agent',
  worker_id: 'agent:mergeos-web-001',
  agent_type: 'frontend-agent',
});

const currentProject = computed(() => {
  if (selectedProjectId.value) {
    return projects.value.find((project) => project.id === selectedProjectId.value) || projects.value[projects.value.length - 1] || null;
  }
  return projects.value[projects.value.length - 1] || null;
});
const currentTasks = computed(() => {
  if (!currentProject.value) return [];
  return tasks.value.filter((task) => task.project_id === currentProject.value.id);
});
const acceptedTasks = computed(() => tasks.value.filter((task) => task.status === 'accepted'));
const openTasks = computed(() => tasks.value.filter((task) => task.status !== 'accepted'));
const recentLedger = computed(() => ledger.value.slice(-8).reverse());
const totalBudget = computed(() => projects.value.reduce((sum, project) => sum + project.budget_cents, 0));
const totalPool = computed(() => projects.value.reduce((sum, project) => sum + project.work_pool_cents, 0));
const tokenSymbol = computed(() => runtimeConfig.value?.token_symbol || 'MERGE');
const statusLabel = computed(() => {
  if (currentProject.value) return `${currentProject.value.payment_provider} verified`;
  return runtimeConfig.value?.repo_provider || 'ready';
});
const paymentReferencePlaceholder = computed(() => {
  if (projectForm.payment_method === 'paypal') return 'PayPal order id';
  if (runtimeConfig.value?.crypto_asset === 'erc20') return 'EVM tx hash for stablecoin transfer';
  return 'EVM tx hash';
});
const budgetUsd = computed({
  get: () => Math.round(projectForm.budget_cents / 100),
  set: (value) => {
    projectForm.budget_cents = Math.max(100, Number(value || 100)) * 100;
  },
});

function money(cents) {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 0,
  }).format((cents || 0) / 100);
}

function shortHash(hash) {
  return `${hash.slice(0, 8)}...${hash.slice(-6)}`;
}

async function api(path, options = {}) {
  const response = await fetch(path, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token.value ? { Authorization: `Bearer ${token.value}` } : {}),
      ...(options.headers || {}),
    },
  });
  const payload = await response.json();
  if (!response.ok) {
    if (response.status === 401) clearSession();
    throw new Error(payload.error || 'Request failed');
  }
  return payload;
}

async function loadConfig() {
  runtimeConfig.value = await api('/api/config');
  if (runtimeConfig.value.dev_payment_enabled && !projectForm.payment_reference) {
    projectForm.payment_reference = runtimeConfig.value.dev_payment_code;
  }
}

async function submitAuth() {
  authBusy.value = true;
  errorMessage.value = '';
  try {
    const path = authMode.value === 'register' ? '/api/auth/register' : '/api/auth/login';
    const body = authMode.value === 'register'
      ? authForm
      : { email: authForm.email, password: authForm.password };
    const auth = await api(path, { method: 'POST', body: JSON.stringify(body) });
    setSession(auth);
    syncProjectContact();
    await refreshProtected();
  } catch (error) {
    errorMessage.value = error.message;
  } finally {
    authBusy.value = false;
  }
}

function setSession(auth) {
  token.value = auth.token;
  user.value = auth.user;
  localStorage.setItem('mergeos_token', auth.token);
}

function clearSession() {
  token.value = '';
  user.value = null;
  projects.value = [];
  tasks.value = [];
  notifications.value = [];
  localStorage.removeItem('mergeos_token');
}

function syncProjectContact() {
  if (!user.value) return;
  projectForm.client_name = user.value.name;
  projectForm.company_name = user.value.company_name;
  projectForm.client_email = user.value.email;
}

async function restoreSession() {
  if (!token.value) return;
  try {
    user.value = await api('/api/auth/me');
    syncProjectContact();
    await refreshProtected();
  } catch {
    clearSession();
  }
}

async function refreshAll() {
  await loadConfig();
  if (user.value) {
    await refreshProtected();
  }
}

async function refreshProtected() {
  const [projectRows, taskRows, ledgerRows, noteRows] = await Promise.all([
    api('/api/projects'),
    api('/api/tasks'),
    api('/api/ledger'),
    api('/api/notifications'),
  ]);
  projects.value = projectRows;
  tasks.value = taskRows;
  ledger.value = ledgerRows;
  notifications.value = noteRows.slice().reverse();
  if (!selectedProjectId.value && projectRows.length) {
    selectedProjectId.value = projectRows[projectRows.length - 1].id;
  }
  reconcileSelection();
}

async function preparePayPalOrder() {
  preparingPayPal.value = true;
  errorMessage.value = '';
  try {
    const order = await api('/api/payments/paypal/orders', {
      method: 'POST',
      body: JSON.stringify({
        amount_cents: projectForm.budget_cents,
        description: projectForm.title,
        return_url: window.location.href,
        cancel_url: window.location.href,
      }),
    });
    paypalOrder.value = order;
    projectForm.payment_reference = order.order_id;
  } catch (error) {
    errorMessage.value = error.message;
  } finally {
    preparingPayPal.value = false;
  }
}

async function createProject() {
  creating.value = true;
  errorMessage.value = '';
  try {
    const project = await api('/api/projects', {
      method: 'POST',
      body: JSON.stringify(projectForm),
    });
    selectedProjectId.value = project.id;
    await refreshProtected();
    portalTab.value = 'workspace';
  } catch (error) {
    errorMessage.value = error.message;
  } finally {
    creating.value = false;
  }
}

function selectProject(project) {
  selectedProjectId.value = project.id;
  reconcileSelection();
}

function selectTask(task) {
  selectedTask.value = task;
  workerForm.worker_kind = task.required_worker_kind;
  workerForm.worker_id = task.required_worker_kind === 'human' ? 'github:client-reviewer' : 'agent:mergeos-web-001';
  workerForm.agent_type = task.required_worker_kind === 'human' ? '' : (task.suggested_agent_type || 'custom-agent');
  errorMessage.value = '';
}

async function acceptSelectedTask() {
  if (!selectedTask.value) return;
  accepting.value = true;
  errorMessage.value = '';
  try {
    await api(`/api/tasks/${selectedTask.value.id}/accept`, {
      method: 'POST',
      body: JSON.stringify(workerForm),
    });
    await refreshProtected();
  } catch (error) {
    errorMessage.value = error.message;
  } finally {
    accepting.value = false;
  }
}

async function logout() {
  try {
    await api('/api/auth/logout', { method: 'POST', body: JSON.stringify({}) });
  } finally {
    clearSession();
  }
}

function reconcileSelection() {
  if (selectedTask.value) {
    const fresh = currentTasks.value.find((task) => task.id === selectedTask.value.id);
    selectedTask.value = fresh || currentTasks.value[0] || null;
    return;
  }
  if (currentTasks.value.length) {
    selectTask(currentTasks.value[0]);
  }
}

watch(currentTasks, reconcileSelection);

onMounted(async () => {
  await loadConfig();
  await restoreSession();
});
</script>
