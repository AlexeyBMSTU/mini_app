
import React from 'react';
import styles from '../styles/common.module.css';

const Loading: React.FC = () => {
  return (
    <div className={styles.loading}>
      <div className={styles.spinner}></div>
      <p>Инициализация приложения...</p>
    </div>
  );
};

export default Loading;