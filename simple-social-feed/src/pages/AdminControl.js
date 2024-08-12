import React, { useState, useEffect } from "react";
import { runMigration, accessProtectedRoute, getMigration } from "../services/api";
import ActionButton from "../components/ActionButton";

function AdminControl() {

    const [warning, setWarning] = useState('')
    const [showConfirm, setShowConfirm] = useState(false)
    const [migrations, setMigrations] = useState([])
    const [selectedMigration, setSelectedMigration] = useState(null)
    const [resultMessage, setResultMessage] = useState('')
    const [isAuthorized, setIsAuthorized] = useState(false)
    const [loading, setLoading] = useState(true)  // This state helps to show a loading message while checking authorization and fetching data from the server

    useEffect(() => {
        const checkAuthorizationAndFetch = async () => {
            try {
                const res = await fetch(accessProtectedRoute)
                if (res.status === 401) {
                    setWarning('Unauthorized')
                } else {
                    setIsAuthorized(true)
                    const migrationRes = await getMigration()
                    setMigrations(migrationRes.data)
                }
            } catch (error) {
                console.error('Error access admin page:', error)
                setWarning('An error occurred')
            } finally {
                setLoading(false)
            }
        }
        checkAuthorizationAndFetch()
    }, [])

    const handleRunMigration = async () => {
        if (!selectedMigration) return;

        setWarning('')
        setShowConfirm(false)

        try {
            const res = await runMigration(selectedMigration)
            setResultMessage((prev) => ({
                ...prev,
                [selectedMigration]: res.data.message
            }))
    } catch (error) {
        console.error('Error running migration:', error)
        setWarning('Failed to run migration. Please try again.')
    } finally {
        setSelectedMigration(null)
    }
}

    const handleConfirmation = (migrationID) => {
        setSelectedMigration(migrationID)
        setShowConfirm(true)
    }

    const handleCancel = () => {
        setShowConfirm(false)
        setSelectedMigration(null)
    }

    if (loading) {
        return <p>Loading...</p>
      }

    if (!isAuthorized) {
        return <p>{warning || "Checking authorization..."}</p>
      }

    return (
        <div>
            <h2>Admin Control</h2>
            { migrations.map((migration) => (
                <div key={migration.migration_id} style={styles.migration}>
                    <p>{migration.migration_title}</p>
                    <ActionButton 
                    text={"Execute"} 
                    onClick={() => handleConfirmation(migration.migration_id)}
                    style={styles.button} 
                    disabled={selectedMigration !== null}
                    />
                    { resultMessage[migration.migration_id] && <p>Successfully execute {resultMessage[migration.migration_id]}</p>}
                </div>
            )) }
            { warning && <p>{warning}</p> }

            { showConfirm && (
                <div style={styles.confirmationModalStyle}>
                    <p>Are you sure you want to run migration?</p>
                    <ActionButton text={'Yes'} onClick={handleRunMigration} style={styles.button}  />
                    <ActionButton text={'Cancel'} onClick={handleCancel} style={styles.button} />
                </div>
            ) }
        </div>
    )
}

const styles = {
    migration: {
        border: "1px solid #ccc",
        padding: "10px",
        marginBottom: "10px",
        borderRadius: "5px",
        position: "relative",
    },
    confirmationModalStyle: {
        backgroundColor: '#f8f8f8',
        border: '1px solid #ccc',
        padding: '20px',
        position: 'fixed',
        top: '50%',
        left: '50%',
        transform: 'translate(-50%, -50%)',
        zIndex: '1000',
        boxShadow: '0 0 10px rgba(0, 0, 0, 0.1)',
    },
    button: {
        margin: '10px',
    }
}

export default AdminControl;