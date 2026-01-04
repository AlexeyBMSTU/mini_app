import { Loader } from '@/components/Loader/Loader'
import { PropertyCard } from '@/components/PropertyCard/PropertyCard'
import { PropertyFilters } from '@/components/PropertyFilters/PropertyFilters'
import { PurePage } from '@/components/PurePage/PurePage'
import { useStore } from '@/store/StoreContext'
import { Box, Typography } from '@mui/material'
import { observer } from 'mobx-react-lite'
import { useEffect } from 'react'
import { ErrorPage } from '../Error/ErrorPage'
import { motion } from 'motion/react'
import styles from '@/pages/Browse/BrowsePage.module.css'

export const BrowsePage = observer(() => {
  const { browse } = useStore()

  useEffect(() => {
    if (browse.properties.length === 0) {
      browse.fetchProperties()
    }
  }, [browse])

  if (browse.loading) {
    return (
      <PurePage>
        <Loader />
      </PurePage>
    )
  }

  if (browse.error) {
    return (
      <PurePage>
        <ErrorPage />
      </PurePage>
    )
  }

  return (
    <PurePage>
      <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} transition={{ duration: 0.5 }}>
        <PropertyFilters />

        {browse.filteredProperties.length === 0 ? (
          <Box className={styles.emptyContainer}>
            <Typography>Не найдено объектов недвижимости по заданным фильтрам</Typography>
          </Box>
        ) : (
          <div className={styles.propertiesGrid}>
            {browse.filteredProperties.map(property => (
              <div className={styles.propertyItem} key={property.id}>
                <PropertyCard property={property} />
              </div>
            ))}
          </div>
        )}
      </motion.div>
    </PurePage>
  )
})
