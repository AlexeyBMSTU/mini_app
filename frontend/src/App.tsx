import {
  DevInstructions,
  Features,
	Footer,
	Header,
	Loading,
	UserInfo,
} from '@/components'
import { useTelegram } from '@/hooks/useTelegram'
import { useApi } from '@/hooks/useApi'
import styles from '@/styles/common.module.css'
import apiService from './services/apiService'
import avitoService from './services/avitoService'

function App() {
	const { user: telegramUser, loading: telegramLoading, isTelegram, showAlert, closeApp } = useTelegram()
	const { loading: apiLoading, error, getUserData, saveUserData } = useApi()

	const loading = telegramLoading || apiLoading

	const handleGetToken = async () => {
		try {
			await avitoService.getItems()
			await avitoService.getChats()
		} catch (error) {
			console.error('Error getting token:', error)
		}
	}

	if (loading) {
		return <Loading />
	}

	return (
		<div className={styles.container}>
			<button className={styles.button} onClick={handleGetToken}>Get Token</button>
		</div>
	)
}

export default App
