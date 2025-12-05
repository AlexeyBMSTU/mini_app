import { useEffect, useState } from 'react';

interface TelegramUser {
  id: number;
  first_name: string;
  last_name?: string;
  username?: string;
  language_code?: string;
  is_premium?: boolean;
}

interface TelegramWebApp {
  ready: () => void;
  expand: () => void;
  close: () => void;
  initDataUnsafe: {
    user?: TelegramUser;
  };
  showAlert: (message: string) => void;
  setHeaderColor: (color: string) => void;
  setBackgroundColor: (color: string) => void;
}

declare global {
  interface Window {
    Telegram?: {
      WebApp: TelegramWebApp;
    };
  }
}

function App() {
  const [user, setUser] = useState<TelegramUser | null>(null);
  const [loading, setLoading] = useState(true);
  const [isTelegram, setIsTelegram] = useState(false);

  useEffect(() => {
    const initApp = () => {
      try {
        const tg = window.Telegram?.WebApp;
        
        const isLocalHost = window.location.hostname === 'localhost' || 
                           window.location.hostname === '127.0.0.1' || 
                           window.location.hostname.startsWith('192.168.');

        if (tg && !isLocalHost) {
          setIsTelegram(true);
          
          tg.ready();
          tg.expand();
          
          const userData = tg.initDataUnsafe?.user;
          if (userData) {
            setUser(userData);
          }
          
          tg.setHeaderColor('#ffffff');
          tg.setBackgroundColor('#ffffff');
          
          console.log('Запущено в Telegram');
        } else {
          setIsTelegram(false);
          console.log('Запущено в браузере. Используется режим разработки.');
          
          setUser({
            id: 123456789,
            first_name: 'Разработчик',
            last_name: 'Тестовый',
            username: 'developer',
            language_code: 'ru',
            is_premium: true
          });
        }
      } catch (error) {
        console.error('Ошибка инициализации:', error);
        setIsTelegram(false);
      } finally {
        setLoading(false);
      }
    };

    initApp();
  }, []);

  const handleTelegramAlert = () => {
    const tg = window.Telegram?.WebApp;
    if (tg) {
      tg.showAlert('Сообщение из Telegram Mini App!');
    } else {
      alert('Сообщение из браузера! (В Telegram будет showAlert)');
    }
  };

  const simulateTelegramAction = () => {
    alert(`Это действие ${isTelegram ? 'в Telegram2' : 'в браузере'}`);
  };

  if (loading) {
    return (
      <div style={styles.loading}>
        <div style={styles.spinner}></div>
        <p>Инициализация приложения...</p>
      </div>
    );
  }

  return (
    <div style={styles.container}>
      <header style={styles.header}>
        <div style={styles.badge}>
          {isTelegram ? '📱 Telegram' : '🖥️ Браузер'}
        </div>
        <h1 style={styles.title}>🚀 Telegram Mini App</h1>
        <p style={styles.subtitle}>
          {isTelegram 
            ? 'Запущено в Telegram2' 
            : 'Режим разработки в браузере'}
        </p>
      </header>
      
      <main style={styles.main}>
        <div style={styles.card}>
          <h2>👤 Информация о пользователе</h2>
          <div style={styles.userInfo}>
            {user ? (
              <>
                <p><strong>ID:</strong> {user.id}</p>
                <p><strong>Имя:</strong> {user.first_name}</p>
                {user.last_name && <p><strong>Фамилия:</strong> {user.last_name}</p>}
                {user.username && <p><strong>Username:</strong> @{user.username}</p>}
                <p><strong>Режим:</strong> {isTelegram ? 'Telegram' : 'Разработка'}</p>
                {!isTelegram && (
                  <p style={styles.devNote}>
                    <small>
                      ⚠️ Это тестовые данные. В Telegram будут реальные данные пользователя.
                    </small>
                  </p>
                )}
              </>
            ) : (
              <p>Данные пользователя не получены</p>
            )}
          </div>
        </div>
        
        <div style={styles.card}>
          <h2>⚡ Демонстрация функций</h2>
          
          <div style={styles.features}>
            <div style={styles.feature}>
              <h3>Telegram API</h3>
              <p>Доступные функции:</p>
              <ul style={styles.featureList}>
                <li>{isTelegram ? '✅' : '⚠️'} Основная кнопка</li>
                <li>{isTelegram ? '✅' : '⚠️'} Telegram Alert</li>
                <li>{isTelegram ? '✅' : '⚠️'} Параметры темы</li>
                <li>✅ Режим разработки</li>
              </ul>
            </div>
            
            <div style={styles.feature}>
              <h3>Технологии</h3>
              <ul style={styles.featureList}>
                <li>✅ React 19</li>
                <li>✅ TypeScript</li>
                <li>✅ Webpack 5</li>
                <li>✅ Telegram Web App</li>
              </ul>
            </div>
          </div>
          
          <div style={styles.buttons}>
            <button 
              style={styles.buttonPrimary}
              onClick={simulateTelegramAction}
            >
              📱 Тестовое действие
            </button>
            
            <button 
              style={styles.buttonSecondary}
              onClick={handleTelegramAlert}
            >
              🔔 Показать уведомление
            </button>
            
            {isTelegram && (
              <button 
                style={styles.buttonTelegram}
                onClick={() => window.Telegram?.WebApp.close()}
              >
                ❌ Закрыть приложение
              </button>
            )}
          </div>
        </div>
        
        {!isTelegram && (
          <div style={styles.devCard}>
            <h3>🛠️ Инструкция для разработки</h3>
            <ol style={styles.instructions}>
              <li>Создайте бота через @BotFather в Telegram</li>
              <li>Настройте Menu Button с URL вашего приложения</li>
              <li>Откройте бота в Telegram и нажмите Menu Button</li>
              <li>Приложение запустится внутри Telegram с реальными данными</li>
            </ol>
            <p style={styles.note}>
              <strong>Примечание:</strong> В режиме разработки используйте моковые данные. 
              В Telegram будут реальные данные пользователя.
            </p>
          </div>
        )}
      </main>
      
      <footer style={styles.footer}>
        <p>
          Telegram Mini App • {isTelegram ? 'Production' : 'Development'} • 
          React + Webpack • {new Date().getFullYear()}
        </p>
      </footer>
    </div>
  );
}

const styles = {
  container: {
    maxWidth: '800px',
    margin: '0 auto',
    padding: '20px',
    minHeight: '100vh',
    display: 'flex',
    flexDirection: 'column' as const,
    fontFamily: 'Arial, sans-serif',
    backgroundColor: '#ffffff',
  },
  loading: {
    display: 'flex',
    flexDirection: 'column' as const,
    justifyContent: 'center',
    alignItems: 'center',
    height: '100vh',
    backgroundColor: '#ffffff',
  },
  spinner: {
    width: '50px',
    height: '50px',
    border: '5px solid #f3f3f3',
    borderTop: '5px solid #3498db',
    borderRadius: '50%',
    animation: 'spin 1s linear infinite',
    marginBottom: '20px',
  },
  header: {
    textAlign: 'center' as const,
    padding: '30px 0',
    borderBottom: '2px solid #f0f0f0',
    marginBottom: '30px',
    position: 'relative' as const,
  },
  badge: {
    position: 'absolute' as const,
    top: '10px',
    right: '10px',
    backgroundColor: '#007bff',
    color: 'white',
    padding: '5px 10px',
    borderRadius: '20px',
    fontSize: '12px',
    fontWeight: 'bold' as const,
  },
  title: {
    fontSize: '2.5rem',
    color: '#333',
    marginBottom: '10px',
  },
  subtitle: {
    fontSize: '1.2rem',
    color: '#666',
  },
  main: {
    flex: 1,
  },
  card: {
    backgroundColor: '#f8f9fa',
    borderRadius: '12px',
    padding: '25px',
    marginBottom: '25px',
    border: '1px solid #e9ecef',
  },
  userInfo: {
    backgroundColor: 'white',
    padding: '20px',
    borderRadius: '8px',
    marginTop: '15px',
  },
  devNote: {
    marginTop: '15px',
    padding: '10px',
    backgroundColor: '#fff3cd',
    border: '1px solid #ffeaa7',
    borderRadius: '5px',
    color: '#856404',
  },
  features: {
    display: 'flex',
    gap: '20px',
    margin: '20px 0',
    flexWrap: 'wrap' as const,
  },
  feature: {
    flex: 1,
    minWidth: '250px',
  },
  featureList: {
    listStyle: 'none',
    padding: '0',
    marginTop: '10px',
  },
  buttons: {
    display: 'flex',
    gap: '15px',
    marginTop: '25px',
    flexWrap: 'wrap' as const,
  },
  buttonPrimary: {
    flex: 1,
    padding: '15px 25px',
    backgroundColor: '#007bff',
    color: 'white',
    border: 'none',
    borderRadius: '8px',
    fontSize: '16px',
    fontWeight: '600' as const,
    cursor: 'pointer',
    minWidth: '200px',
    transition: 'all 0.2s ease',
  },
  buttonSecondary: {
    flex: 1,
    padding: '15px 25px',
    backgroundColor: '#6c757d',
    color: 'white',
    border: 'none',
    borderRadius: '8px',
    fontSize: '16px',
    fontWeight: '600' as const,
    cursor: 'pointer',
    minWidth: '200px',
    transition: 'all 0.2s ease',
  },
  buttonTelegram: {
    flex: 1,
    padding: '15px 25px',
    backgroundColor: '#3399ff',
    color: 'white',
    border: 'none',
    borderRadius: '8px',
    fontSize: '16px',
    fontWeight: '600' as const,
    cursor: 'pointer',
    minWidth: '200px',
    transition: 'all 0.2s ease',
  },
  devCard: {
    backgroundColor: '#e7f5ff',
    borderRadius: '12px',
    padding: '25px',
    marginTop: '30px',
    borderLeft: '5px solid #3399ff',
  },
  instructions: {
    margin: '15px 0',
    paddingLeft: '20px',
  },
  note: {
    marginTop: '15px',
    padding: '10px',
    backgroundColor: '#d1ecf1',
    border: '1px solid #bee5eb',
    borderRadius: '5px',
    color: '#0c5460',
  },
  footer: {
    textAlign: 'center' as const,
    padding: '30px 0',
    marginTop: '40px',
    borderTop: '2px solid #f0f0f0',
    color: '#6c757d',
    fontSize: '14px',
  },
};

const style = document.createElement('style');
style.textContent = `
  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
  
  button:hover {
    opacity: 0.9;
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  }
  
  @media (max-width: 768px) {
    .buttons, .features {
      flex-direction: column;
    }
    
    button, .feature {
      width: 100%;
      min-width: auto;
    }
    
    .title {
      font-size: 2rem;
    }
  }
  
  li {
    margin: 8px 0;
    padding-left: 5px;
  }
`;
document.head.appendChild(style);

export default App;