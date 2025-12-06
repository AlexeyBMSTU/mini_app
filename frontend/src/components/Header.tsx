

import React from 'react';
import styles from '../styles/common.module.css';

interface HeaderProps {
  isTelegram: boolean;
}

const Header: React.FC<HeaderProps> = ({ isTelegram }) => {
  return (
    <header className={styles.header}>
      <div className={styles.badge}>
        {isTelegram ? '📱 Telegram' : '🖥️ Браузер'}
      </div>
      <h1 className={styles.title}>🚀 Telegram Mini App</h1>
      <p className={styles.subtitle}>
        {isTelegram
          ? 'Запущено в Telegram2'
          : 'Режим разработки в браузере'}
      </p>
    </header>
  );
};

export default Header;