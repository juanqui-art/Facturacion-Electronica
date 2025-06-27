// Tipos globales para el sistema de facturación SRI

// Extensiones del objeto Window
declare global {
  interface Window {
    // Toast Manager
    toastManager: {
      show: (message: string, type?: 'success' | 'error' | 'warning' | 'info', duration?: number, options?: ToastOptions) => string | null;
      hide: (id: string) => void;
      success: (message: string, options?: ToastOptions) => string | null;
      error: (message: string, options?: ToastOptions) => string | null;
      warning: (message: string, options?: ToastOptions) => string | null;
      info: (message: string, options?: ToastOptions) => string | null;
      clear: () => void;
    };
    
    // Funciones globales de conveniencia
    showToast: (message: string, type?: 'success' | 'error' | 'warning' | 'info', duration?: number, options?: ToastOptions) => string | null;
    showSuccess: (message: string, options?: ToastOptions) => string | null;
    showError: (message: string, options?: ToastOptions) => string | null;
    showWarning: (message: string, options?: ToastOptions) => string | null;
    showInfo: (message: string, options?: ToastOptions) => string | null;
  }
}

// Tipos para el sistema de Toast
export interface ToastOptions {
  title?: string;
  persistent?: boolean;
  action?: string | (() => void);
  actionText?: string;
}

// Tipos para elementos DOM comunes
export type HTMLInputElementWithValue = HTMLInputElement | null;
export type HTMLSelectElementWithValue = HTMLSelectElement | null;
export type HTMLTextAreaElementWithValue = HTMLTextAreaElement | null;

// Tipos para eventos comunes
export type InputChangeEvent = Event & {
  target: HTMLInputElement;
};

export type SelectChangeEvent = Event & {
  target: HTMLSelectElement;
};

// Tipos para datos de la API
export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  message?: string;
  error?: string;
}

export interface Cliente {
  id: number;
  cedula: string;
  nombre: string;
  direccion?: string;
  telefono?: string;
  email?: string;
  tipoCliente: 'PERSONA_NATURAL' | 'EMPRESA';
  fechaCreacion: string;
  activo: boolean;
}

export interface Factura {
  id: number;
  numeroFactura: string;
  claveAcceso: string;
  fechaEmision: string;
  clienteNombre: string;
  clienteCedula: string;
  subtotal: number;
  iva: number;
  total: number;
  estado: 'BORRADOR' | 'ENVIADA' | 'AUTORIZADA' | 'RECHAZADA';
  numeroAutorizacion?: string;
  ambiente: 'PRUEBAS' | 'PRODUCCION';
}

export interface Estadisticas {
  total_facturas: number;
  por_estado: {
    BORRADOR?: number;
    ENVIADA?: number;
    AUTORIZADA?: number;
    RECHAZADA?: number;
  };
  total_facturado: number;
}

export interface EstadoSRI {
  disponible: boolean;
  mensaje: string;
  ambiente?: string;
  timestamp: string;
}

// Helper type para elementos del DOM que pueden ser null
export type ElementOrNull<T extends Element> = T | null;

// Para uso sin export (hacer disponible globalmente)
export {};