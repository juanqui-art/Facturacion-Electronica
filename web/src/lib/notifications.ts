/**
 * Sistema de notificaciones mejorado para feedback al usuario
 * Integrado con el sistema de componentes Astro + Tailwind
 */

export type NotificationType = 'success' | 'error' | 'warning' | 'info';

export interface Notification {
  id: string;
  type: NotificationType;
  title: string;
  message?: string;
  duration?: number;
  autoClose?: boolean;
}

class NotificationManager {
  private notifications: Map<string, Notification> = new Map();
  private callbacks: Set<(notifications: Notification[]) => void> = new Set();
  private idCounter = 0;

  private generateId(): string {
    return `notification-${++this.idCounter}-${Date.now()}`;
  }

  show(notification: Omit<Notification, 'id'>): string {
    const id = this.generateId();
    const fullNotification: Notification = {
      id,
      autoClose: true,
      duration: 5000,
      ...notification
    };

    this.notifications.set(id, fullNotification);
    this.notifyCallbacks();

    // Auto-close si está habilitado
    if (fullNotification.autoClose && fullNotification.duration) {
      setTimeout(() => {
        this.remove(id);
      }, fullNotification.duration);
    }

    return id;
  }

  remove(id: string): void {
    if (this.notifications.delete(id)) {
      this.notifyCallbacks();
    }
  }

  clear(): void {
    this.notifications.clear();
    this.notifyCallbacks();
  }

  getAll(): Notification[] {
    return Array.from(this.notifications.values());
  }

  subscribe(callback: (notifications: Notification[]) => void): () => void {
    this.callbacks.add(callback);
    return () => this.callbacks.delete(callback);
  }

  private notifyCallbacks(): void {
    const notifications = this.getAll();
    this.callbacks.forEach(callback => callback(notifications));
  }

  // Métodos de conveniencia
  success(title: string, message?: string): string {
    return this.show({ type: 'success', title, message });
  }

  error(title: string, message?: string): string {
    return this.show({ 
      type: 'error', 
      title, 
      message, 
      duration: 8000 // Errores duran más
    });
  }

  warning(title: string, message?: string): string {
    return this.show({ type: 'warning', title, message });
  }

  info(title: string, message?: string): string {
    return this.show({ type: 'info', title, message });
  }

  // Para errores de API específicos
  apiError(operation: string, error: string): string {
    return this.error(
      `Error en ${operation}`,
      `No se pudo completar la operación: ${error}`
    );
  }

  // Para operaciones exitosas de API
  apiSuccess(operation: string, details?: string): string {
    return this.success(
      `${operation} exitoso`,
      details
    );
  }
}

// Instancia global del manager
export const notifications = new NotificationManager();

// Helper para manejar respuestas de API automáticamente
export function handleApiResponse<T>(
  response: import('./api').ApiResponse<T>,
  successMessage?: string,
  operation?: string
): T | null {
  if (response.success && response.data) {
    if (successMessage) {
      notifications.success(successMessage);
    }
    return response.data;
  } else {
    if (operation && response.error) {
      notifications.apiError(operation, response.error);
    }
    return null;
  }
}

// CSS classes para los diferentes tipos de notificación
export const notificationStyles = {
  success: 'bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-800 text-green-800 dark:text-green-200',
  error: 'bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800 text-red-800 dark:text-red-200',
  warning: 'bg-amber-50 dark:bg-amber-900/20 border-amber-200 dark:border-amber-800 text-amber-800 dark:text-amber-200',
  info: 'bg-blue-50 dark:bg-blue-900/20 border-blue-200 dark:border-blue-800 text-blue-800 dark:text-blue-200'
};

// Icons para cada tipo
export const notificationIcons = {
  success: `<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
  </svg>`,
  error: `<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
  </svg>`,
  warning: `<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
    <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
  </svg>`,
  info: `<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
    <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
  </svg>`
};