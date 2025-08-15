package admin

import (
	"dailybot/internal/config"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type SimpleAdmin struct {
	config    *config.Config
	stats     Stats
	mu        sync.RWMutex
	startTime time.Time
}

type Stats struct {
	TotalMessages    int64
	WeatherRequests  int64
	NewsRequests     int64
	ExchangeRequests int64
	ActiveUsers      map[int64]time.Time
}

func NewSimpleAdmin(cfg *config.Config) *SimpleAdmin {
	return &SimpleAdmin{
		config:    cfg,
		startTime: time.Now(),
		stats: Stats{
			ActiveUsers: make(map[int64]time.Time),
		},
	}
}

func (a *SimpleAdmin) Start() {
	http.HandleFunc("/", a.handleAdmin)
	http.HandleFunc("/api/stats", a.handleStats)

	port := a.config.AdminPort
	log.Printf("🎛 Admin panel: http://localhost:%s", port)
	log.Printf("🔑 Password: %s", a.config.AdminPassword)

	http.ListenAndServe(":"+port, nil)
}

func (a *SimpleAdmin) LogCommand(userID int64, command, args string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.stats.TotalMessages++
	a.stats.ActiveUsers[userID] = time.Now()

	switch command {
	case "weather":
		a.stats.WeatherRequests++
	case "news":
		a.stats.NewsRequests++
	case "exchange":
		a.stats.ExchangeRequests++
	}
}

func (a *SimpleAdmin) handleAdmin(w http.ResponseWriter, r *http.Request) {
	// Простая авторизация через форму
	if r.Method == "POST" {
		password := r.FormValue("password")
		if password == a.config.AdminPassword {
			// Устанавливаем куку на 24 часа
			http.SetCookie(w, &http.Cookie{
				Name:     "admin_auth",
				Value:    a.config.AdminPassword,
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true,
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	// Проверяем авторизацию
	cookie, err := r.Cookie("admin_auth")
	if err != nil || cookie.Value != a.config.AdminPassword {
		a.showLoginForm(w)
		return
	}

	// Проверяем logout
	if r.URL.Query().Get("logout") == "1" {
		http.SetCookie(w, &http.Cookie{
			Name:     "admin_auth",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			HttpOnly: true,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	a.showDashboard(w)
}

func (a *SimpleAdmin) showLoginForm(w http.ResponseWriter) {
	html := `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>DailyBot Admin - Вход</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .login-card {
            background: rgba(255, 255, 255, 0.95);
            border-radius: 20px;
            backdrop-filter: blur(10px);
            box-shadow: 0 15px 35px rgba(0, 0, 0, 0.1);
            padding: 40px;
            max-width: 400px;
            width: 90%;
            text-align: center;
        }
        .logo { font-size: 3rem; margin-bottom: 10px; }
        h1 { color: #333; margin-bottom: 30px; font-size: 1.8rem; }
        input {
            width: 100%;
            padding: 15px;
            margin: 15px 0;
            border: 2px solid #e2e8f0;
            border-radius: 10px;
            font-size: 16px;
            transition: border-color 0.3s ease;
        }
        input:focus {
            outline: none;
            border-color: #667eea;
        }
        button {
            width: 100%;
            padding: 15px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            border-radius: 10px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: transform 0.3s ease;
        }
        button:hover { transform: translateY(-2px); }
        .info { margin-top: 20px; color: #666; font-size: 14px; }
    </style>
</head>
<body>
    <div class="login-card">
        <div class="logo">🤖</div>
        <h1>DailyBot Admin</h1>
        <form method="POST">
            <input type="password" name="password" placeholder="Введите пароль администратора" required>
            <button type="submit">Войти в панель</button>
        </form>
        <div class="info">
            Панель управления Telegram-ботом<br>
            Только для администраторов
        </div>
    </div>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func (a *SimpleAdmin) showDashboard(w http.ResponseWriter) {
	a.mu.RLock()
	stats := a.stats
	a.mu.RUnlock()

	// Считаем активных пользователей за 24 часа
	activeCount := 0
	now := time.Now()
	for _, lastSeen := range stats.ActiveUsers {
		if now.Sub(lastSeen) < 24*time.Hour {
			activeCount++
		}
	}

	uptime := time.Since(a.startTime)

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>DailyBot Admin - Панель управления</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        .container { max-width: 1200px; margin: 0 auto; }
        .header {
            background: rgba(255, 255, 255, 0.1);
            backdrop-filter: blur(10px);
            color: white;
            padding: 20px;
            border-radius: 15px;
            margin-bottom: 20px;
            display: flex;
            justify-content: space-between;
            align-items: center;
            flex-wrap: wrap;
        }
        .header h1 { font-size: 2rem; }
        .header .actions { display: flex; gap: 10px; }
        .btn {
            background: rgba(255, 255, 255, 0.2);
            color: white;
            padding: 10px 20px;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            text-decoration: none;
            font-weight: 500;
            transition: background 0.3s ease;
        }
        .btn:hover { background: rgba(255, 255, 255, 0.3); }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
            gap: 20px;
            margin-bottom: 20px;
        }
        .stat-card {
            background: rgba(255, 255, 255, 0.95);
            border-radius: 15px;
            padding: 25px;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
            backdrop-filter: blur(10px);
            transition: transform 0.3s ease;
        }
        .stat-card:hover { transform: translateY(-5px); }
        .stat-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 15px;
        }
        .stat-icon {
            font-size: 2.5rem;
            opacity: 0.7;
        }
        .stat-number {
            font-size: 2.5rem;
            font-weight: bold;
            color: #333;
            margin-bottom: 5px;
        }
        .stat-label {
            color: #666;
            font-size: 0.95rem;
        }
        .info-card {
            background: rgba(255, 255, 255, 0.95);
            border-radius: 15px;
            padding: 25px;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
            backdrop-filter: blur(10px);
        }
        .status { color: #22c55e; font-weight: 600; }
        .status::before { content: "🟢 "; }
        .auto-refresh {
            background: #f0f9ff;
            border: 1px solid #0ea5e9;
            border-radius: 8px;
            padding: 12px;
            margin-top: 20px;
            color: #0369a1;
            font-size: 14px;
            text-align: center;
        }
        @media (max-width: 768px) {
            .header { flex-direction: column; gap: 15px; text-align: center; }
            .stats-grid { grid-template-columns: 1fr; }
            body { padding: 10px; }
        }
    </style>
    <script>
        // Автообновление каждые 30 секунд
        setTimeout(() => location.reload(), 30000);
        
        // Показываем время до следующего обновления
        let countdown = 30;
        setInterval(() => {
            countdown--;
            const element = document.getElementById('countdown');
            if (element) element.textContent = countdown;
            if (countdown <= 0) countdown = 30;
        }, 1000);
    </script>
</head>
<body>
    <div class="container">
        <div class="header">
            <div>
                <h1>🤖 DailyBot Admin</h1>
                <p>Панель управления Telegram-ботом</p>
            </div>
            <div class="actions">
                <button class="btn" onclick="location.reload()">🔄 Обновить</button>
                <a href="?logout=1" class="btn">🚪 Выход</a>
            </div>
        </div>
        
        <div class="stats-grid">
            <div class="stat-card">
                <div class="stat-header">
                    <div>
                        <div class="stat-number">%d</div>
                        <div class="stat-label">💬 Всего сообщений</div>
                    </div>
                    <div class="stat-icon">💬</div>
                </div>
            </div>
            
            <div class="stat-card">
                <div class="stat-header">
                    <div>
                        <div class="stat-number">%d</div>
                        <div class="stat-label">👥 Активных пользователей</div>
                    </div>
                    <div class="stat-icon">👥</div>
                </div>
                <small style="color: #666;">из %d всего (за 24ч)</small>
            </div>
            
            <div class="stat-card">
                <div class="stat-header">
                    <div>
                        <div class="stat-number">%d</div>
                        <div class="stat-label">🌤 Запросов погоды</div>
                    </div>
                    <div class="stat-icon">🌤</div>
                </div>
            </div>
            
            <div class="stat-card">
                <div class="stat-header">
                    <div>
                        <div class="stat-number">%d</div>
                        <div class="stat-label">📰 Запросов новостей</div>
                    </div>
                    <div class="stat-icon">📰</div>
                </div>
            </div>
            
            <div class="stat-card">
                <div class="stat-header">
                    <div>
                        <div class="stat-number">%d</div>
                        <div class="stat-label">💱 Запросов валют</div>
                    </div>
                    <div class="stat-icon">💱</div>
                </div>
            </div>
            
            <div class="stat-card">
                <div class="stat-header">
                    <div>
                        <div class="stat-number">%s</div>
                        <div class="stat-label">⏱ Время работы</div>
                    </div>
                    <div class="stat-icon">⏱</div>
                </div>
            </div>
        </div>
        
        <div class="info-card">
            <h3 style="margin-bottom: 20px;">📊 Информация о боте</h3>
            <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px;">
                <div>
                    <strong>Статус:</strong><br>
                    <span class="status">Онлайн и работает</span>
                </div>
                <div>
                    <strong>Запущен:</strong><br>
                    %s
                </div>
                <div>
                    <strong>Версия:</strong><br>
                    DailyBot v1.0.0
                </div>
                <div>
                    <strong>Админка:</strong><br>
                    Simple Web Interface
                </div>
            </div>
            
            <div class="auto-refresh">
                🔄 Автоматическое обновление через <span id="countdown">30</span> секунд
            </div>
        </div>
    </div>
</body>
</html>`,
		stats.TotalMessages,
		activeCount,
		len(stats.ActiveUsers),
		stats.WeatherRequests,
		stats.NewsRequests,
		stats.ExchangeRequests,
		formatDuration(uptime),
		a.startTime.Format("02.01.2006 15:04:05"),
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func (a *SimpleAdmin) handleStats(w http.ResponseWriter, r *http.Request) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	activeCount := 0
	now := time.Now()
	for _, lastSeen := range a.stats.ActiveUsers {
		if now.Sub(lastSeen) < 24*time.Hour {
			activeCount++
		}
	}

	uptime := time.Since(a.startTime)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Fprintf(w, `{
		"status": "ok",
		"totalMessages": %d,
		"activeUsers": %d,
		"totalUsers": %d,
		"weatherRequests": %d,
		"newsRequests": %d,
		"exchangeRequests": %d,
		"uptimeSeconds": %.0f,
		"uptimeFormatted": "%s",
		"startTime": "%s"
	}`, a.stats.TotalMessages, activeCount, len(a.stats.ActiveUsers),
		a.stats.WeatherRequests, a.stats.NewsRequests, a.stats.ExchangeRequests,
		uptime.Seconds(), formatDuration(uptime), a.startTime.Format("2006-01-02 15:04:05"))
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dд %dч %dм", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dч %dм", hours, minutes)
	}
	return fmt.Sprintf("%dм", minutes)
}
