
import React from 'react';
import cardStyles from '../styles/card.module.css';
import buttonStyles from '../styles/button.module.css';

interface FeaturesProps {
  isTelegram: boolean;
  onTestAction: () => void;
  onShowAlert: () => void;
  onCloseApp: () => void;
  onGetUserData?: () => void;
  onSaveUserData?: () => void;
}

const Features: React.FC<FeaturesProps> = ({
  isTelegram,
  onTestAction,
  onShowAlert,
  onCloseApp,
  onGetUserData,
  onSaveUserData
}) => {
  return (
    <div className={cardStyles.card}>
      <h2>⚡ Демонстрация функций</h2>
      
      <div className={cardStyles.features}>
        <div className={cardStyles.feature}>
          <h3>Telegram API</h3>
          <p>Доступные функции:</p>
          <ul className={cardStyles.featureList}>
            <li>{isTelegram ? '✅' : '⚠️'} Основная кнопка</li>
            <li>{isTelegram ? '✅' : '⚠️'} Telegram Alert</li>
            <li>{isTelegram ? '✅' : '⚠️'} Параметры темы</li>
            <li>✅ Режим разработки</li>
          </ul>
        </div>
        
        <div className={cardStyles.feature}>
          <h3>Технологии</h3>
          <ul className={cardStyles.featureList}>
            <li>✅ React 19</li>
            <li>✅ TypeScript</li>
            <li>✅ Webpack 5</li>
            <li>✅ Telegram Web App</li>
          </ul>
        </div>
      </div>
      
      <div className={buttonStyles.buttons}>
        <button
          className={buttonStyles.buttonPrimary}
          onClick={onTestAction}
        >
          📱 Тестовое действие
        </button>
        
        <button
          className={buttonStyles.buttonSecondary}
          onClick={onShowAlert}
        >
          🔔 Показать уведомление
        </button>
        
        {isTelegram && (
          <button
            className={buttonStyles.buttonTelegram}
            onClick={onCloseApp}
          >
            ❌ Закрыть приложение
          </button>
        )}
        
        {onGetUserData && (
          <button
            className={buttonStyles.buttonSecondary}
            onClick={onGetUserData}
          >
            📥 Получить данные
          </button>
        )}
        
        {onSaveUserData && (
          <button
            className={buttonStyles.buttonSecondary}
            onClick={onSaveUserData}
          >
            💾 Сохранить данные
          </button>
        )}
      </div>
    </div>
  );
};

export default Features;