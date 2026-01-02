import { FormGroup, FormLabel, TextField } from '@mui/material'
import clsx from 'clsx'
import styles from './Form.module.css'

export const Form = ({ label, className }: { label?: string; className?: string }) => {
  return (
    <FormGroup className={clsx(styles.root, className)}>
      {label && <FormLabel className={styles.label}>{label}</FormLabel>}
      <TextField className={styles.input} fullWidth label='ClientID' variant='outlined' />
      <TextField className={styles.input} fullWidth label='ClientSecret' variant='outlined' />
    </FormGroup>
  )
}
