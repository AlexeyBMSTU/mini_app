import { motion } from 'framer-motion'

export const TypingEffect = ({ text, isBold }: { text: string; isBold?: boolean }) => {
  return (
    <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} transition={{ duration: 0.5 }}>
      {text.split('').map((char, index) => (
        <motion.span
          key={index}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{
            duration: 0.1,
            delay: index * 0.05,
          }}
          style={{ fontWeight: isBold ? 'bold' : '' }}
        >
          {char}
        </motion.span>
      ))}
    </motion.div>
  )
}
