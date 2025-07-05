import Modal from './Modal';
import '../../styles/common/ConfirmDialog.css';

const ConfirmDialog = ({ 
  isOpen, 
  title, 
  message, 
  onConfirm, 
  onCancel, 
  confirmText = "Confirmar", 
  cancelText = "Cancelar",
  type = "warning",
  disabled = false 
}) => {
  return (
    <Modal isOpen={isOpen} onClose={onCancel} className="confirm-dialog">
      <div className="confirm-content">
        <div className={`confirm-icon ${type}`}>
          {type === 'danger' && (
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
            </svg>
          )}
          {type === 'warning' && (
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z"/>
            </svg>
          )}
          {type === 'info' && (
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-6h2v6zm0-8h-2V7h2v2z"/>
            </svg>
          )}
        </div>
        
        <h3 className="confirm-title">{title}</h3>
        
        <div className="confirm-message">
          {typeof message === 'string' ? <p>{message}</p> : message}
        </div>
        
        <div className="confirm-actions">
          <button 
            className={`btn btn-${type}`}
            onClick={onConfirm}
            disabled={disabled}
          >
            {confirmText}
          </button>
          {cancelText && (
            <button 
              className="btn btn-secondary"
              onClick={onCancel}
              disabled={disabled}
            >
              {cancelText}
            </button>
          )}
        </div>
      </div>
    </Modal>
  );
};

export default ConfirmDialog; 