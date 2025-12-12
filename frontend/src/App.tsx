import {
  DevInstructions,
  Features,
	Footer,
	Header,
	Loading,
	UserInfo,
} from '@/components'
import { useTelegram } from '@/hooks/useTelegram'
import styles from '@/styles/common.module.css'

function App() {
	const { user, loading, isTelegram, showAlert, closeApp } = useTelegram()

	const handleTelegramAlert = () => {
		showAlert('Сообщение из Telegram Mini App!')
	}

	const simulateTelegramAction = () => {
		alert(`Это действие ${isTelegram ? 'в Telegram2' : 'в браузере'}`)
	}

	if (loading) {
		return <Loading />
	}

	return (
		<div className={styles.container}>
			<Header isTelegram={isTelegram} />

			<main className={styles.main}>
				<UserInfo user={user} isTelegram={isTelegram} />

				<Features
					isTelegram={isTelegram}
					onTestAction={simulateTelegramAction}
					onShowAlert={handleTelegramAlert}
					onCloseApp={closeApp}
				/>

				{!isTelegram && <DevInstructions />}
			</main>

			<Footer isTelegram={isTelegram} />
		</div>
	)
}

export default App
