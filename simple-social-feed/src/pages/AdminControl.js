import React, { useState, useEffect } from "react";
import { adminMigration, accessProtectedRoute } from "../services/api";
import ActionButton from "../components/ActionButton";

function AdminControl() {

    const [warning, setWarning] = useState('')
    const [showConfirm, setShowConfirm] = useState(false)
    const [isAuthorized, setIsAuthorized] = useState(false)
    const [loading, setLoading] = useState(true)  // This state helps to show a loading message while checking authorization and fetching data from the server

    useEffect(() => {
        const checkAuthorization = async () => {
            try {
                const res = await fetch(accessProtectedRoute)
                if (res.status === 401) {
                    setWarning('Unauthorized')
                } else if (res.ok) {
                    setIsAuthorized(true)
                }
            } catch (error) {
                console.error('Error access admin page:', error)
                setWarning('An error occurred')
            } finally {
                setLoading(false)
            }
        }
        checkAuthorization()
    }, [])

    const handleRunMigration = async () => {
        setWarning('')

        try {
            const res = await adminMigration()
            setWarning(res.data.message)
            setShowConfirm(false)
        } catch (error) {
            setWarning(error.response.data.message)
            setShowConfirm(false)
        }
    }

    const handleConfirmation = () => {
        setShowConfirm(true)
    }

    const handleCancel = () => {
        setShowConfirm(false)
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
            <ActionButton text={'Run Migrations'} onClick={handleConfirmation} />
            { showConfirm && (
                <div style={styles.confirmationModalStyle}>
                    <p>Are you sure you want to run migration?</p>
                    <ActionButton text={'Yes'} onClick={handleRunMigration} style={styles.button}  />
                    <ActionButton text={'Cancel'} onClick={handleCancel} style={styles.button} />
                </div>
            ) }
            { warning && <p>{warning}</p> }           
        </div>
    )
}

const styles = {
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