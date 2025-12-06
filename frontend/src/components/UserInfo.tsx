import React from 'react';
import { TelegramUser } from '../types/telegram';
import cardStyles from '../styles/card.module.css';

interface UserInfoProps {
  user: TelegramUser | null;
  isTelegram: boolean;
}

const UserInfo: React.FC<UserInfoProps> = ({ user, isTelegram }) => {
  return (
    <div className={cardStyles.card}>
      <h2>👤 Информация о пользователе</h2>
      <div className={cardStyles.userInfo}>
        {user ? (
          <>
            <p><strong>ID:</strong> {user.id}</p>
            <p><strong>Имя:</strong> {user.first_name}</p>
            {user.last_name && <p><strong>Фамилия:</strong> {user.last_name}</p>}
            {user.username && <p><strong>Username:</strong> @{user.username}</p>}
            <p><strong>Режим:</strong> {isTelegram ? 'Telegram' : 'Разработка'}</p>
            {!isTelegram && (
              <p className={cardStyles.devNote}>
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
  );
};

export default UserInfo;