
import React from 'react';
import cardStyles from '../styles/card.module.css';

const DevInstructions: React.FC = () => {
  return (
    <div className={cardStyles.devCard}>
      <h3>🛠️ Инструкция для разработки</h3>
      <ol className={cardStyles.instructions}>
        <li>Создайте бота через @BotFather в Telegram</li>
        <li>Настройте Menu Button с URL вашего приложения</li>
        <li>Откройте бота в Telegram и нажмите Menu Button</li>
        <li>Приложение запустится внутри Telegram с реальными данными</li>
      </ol>
      <p className={cardStyles.note}>
        <strong>Примечание:</strong> В режиме разработки используйте моковые данные.
        В Telegram будут реальные данные пользователя.
      </p>
    </div>
  );
};

export default DevInstructions;