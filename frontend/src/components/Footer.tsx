
import React from 'react';
import styles from '../styles/common.module.css';

interface FooterProps {
  isTelegram: boolean;
}

const Footer: React.FC<FooterProps> = ({ isTelegram }) => {
  return (
    <footer className={styles.footer}>
      <p>
        Telegram Mini App • {isTelegram ? 'Production' : 'Development'} •
        React + Webpack • {new Date().getFullYear()}
      </p>
    </footer>
  );
};

export default Footer;