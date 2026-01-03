import { Loader } from '@/components/Loader/Loader'
import { PurePage } from '@/components/PurePage/PurePage'
import { useStore } from '@/store/StoreContext'
import { observer } from 'mobx-react-lite'
import { useEffect } from 'react'
import { ErrorPage } from '../Error/ErrorPage'
import styles from './BrowsePage.module.css'
import { PropertyCard } from '@/components/PropertyCard/PropertyCard'
import { PropertyFilters } from '@/components/PropertyFilters/PropertyFilters'
import { Box, Typography } from '@mui/material'

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
      <Box className={styles.browseContainer}>        
        <PropertyFilters />
        
        {browse.filteredProperties.length === 0 ? (
          <Box className={styles.emptyContainer}>
            <Typography>Не найдено объектов недвижимости по заданным фильтрам</Typography>
          </Box>
        ) : (
          <div className={styles.propertiesGrid}>
            {browse.filteredProperties.map((property) => (
              <div className={styles.propertyItem} key={property.id}>
                <PropertyCard property={property} />
              </div>
            ))}
          </div>
        )}
      </Box>
    </PurePage>
  )
})
