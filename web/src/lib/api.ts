/**
 * Cliente API centralizado para comunicación con backend Go
 * Optimizado para el stack Go + Astro
 */

// Configuración de API inteligente
const getApiBase = (): string => {
  if (typeof window === 'undefined') return ''; // SSR
  
  // Desarrollo: Astro en 4321, Go en 8080
  if (window.location.port === '4321') {
    return 'http://localhost:8080/api';
  }
  
  // Producción: Todo en el mismo puerto
  return window.location.origin + '/api';
};

// Tipos para respuestas del API
export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

export interface EstadisticasData {
  total_facturas: number;
  autorizadas: number;
  pendientes: number;
  rechazadas: number;
}

export interface FacturaData {
  id?: number;
  numeroFactura: string;
  clienteNombre: string;
  total: string;
  estado: 'BORRADOR' | 'ENVIADA' | 'AUTORIZADA' | 'RECHAZADA';
  fechaCreacion?: string;
}

export interface SriStatusData {
  disponible: boolean;
  mensaje: string;
  timestamp: string;
  ambiente?: string;
}

// Cliente API mejorado con manejo de errores y cache básico
class ApiClient {
  private baseUrl: string;
  private cache: Map<string, { data: any; timestamp: number }> = new Map();
  private readonly CACHE_TTL = 5 * 60 * 1000; // 5 minutos

  constructor() {
    this.baseUrl = getApiBase();
  }

  // Método base para hacer requests
  private async request<T>(
    endpoint: string, 
    options: RequestInit = {},
    useCache = false
  ): Promise<ApiResponse<T>> {
    try {
      // Verificar cache
      if (useCache) {
        const cached = this.getFromCache(endpoint);
        if (cached) return cached;
      }

      const url = `${this.baseUrl}${endpoint}`;
      const response = await fetch(url, {
        headers: {
          'Content-Type': 'application/json',
          ...options.headers,
        },
        ...options,
      });

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      const data: ApiResponse<T> = await response.json();
      
      // Guardar en cache si es exitoso
      if (useCache && data.success) {
        this.setCache(endpoint, data);
      }

      return data;
    } catch (error) {
      console.error(`API Error (${endpoint}):`, error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Error desconocido'
      };
    }
  }

  // Cache helpers
  private getFromCache<T>(key: string): ApiResponse<T> | null {
    const cached = this.cache.get(key);
    if (!cached) return null;

    if (Date.now() - cached.timestamp > this.CACHE_TTL) {
      this.cache.delete(key);
      return null;
    }

    return cached.data;
  }

  private setCache<T>(key: string, data: ApiResponse<T>): void {
    this.cache.set(key, {
      data,
      timestamp: Date.now()
    });
  }

  // Métodos específicos del API
  async getEstadisticas(): Promise<ApiResponse<EstadisticasData>> {
    return this.request<EstadisticasData>('/estadisticas', {}, true);
  }

  async getFacturas(limit?: number): Promise<ApiResponse<{ facturas: FacturaData[] }>> {
    const endpoint = limit ? `/facturas/db/list?limit=${limit}` : '/facturas/db/list';
    return this.request<{ facturas: FacturaData[] }>(endpoint, {}, true);
  }

  async getFactura(id: string): Promise<ApiResponse<FacturaData>> {
    return this.request<FacturaData>(`/facturas/db/${id}`);
  }

  async createFactura(factura: Partial<FacturaData>): Promise<ApiResponse<FacturaData>> {
    return this.request<FacturaData>('/facturas/db', {
      method: 'POST',
      body: JSON.stringify(factura)
    });
  }

  async updateFactura(id: string, factura: Partial<FacturaData>): Promise<ApiResponse<FacturaData>> {
    return this.request<FacturaData>(`/facturas/db/${id}`, {
      method: 'PUT',
      body: JSON.stringify(factura)
    });
  }

  async deleteFactura(id: string): Promise<ApiResponse<void>> {
    return this.request<void>(`/facturas/db/${id}`, {
      method: 'DELETE'
    });
  }

  async getSriStatus(): Promise<ApiResponse<SriStatusData>> {
    return this.request<SriStatusData>('/sri/status');
  }

  // Limpiar cache manualmente
  clearCache(): void {
    this.cache.clear();
  }

  // Health check
  async healthCheck(): Promise<boolean> {
    try {
      const response = await fetch(`${this.baseUrl}/health`);
      return response.ok;
    } catch {
      return false;
    }
  }
}

// Instancia singleton del cliente
export const api = new ApiClient();

// Hook-like helper para manejo de estado en el frontend
export class DataManager<T> {
  private data: T | null = null;
  private loading = false;
  private error: string | null = null;
  private callbacks: Set<() => void> = new Set();

  constructor(private fetchFn: () => Promise<ApiResponse<T>>) {}

  async load(): Promise<void> {
    this.setLoading(true);
    this.setError(null);

    const result = await this.fetchFn();
    
    if (result.success && result.data) {
      this.setData(result.data);
    } else {
      this.setError(result.error || 'Error al cargar datos');
    }
    
    this.setLoading(false);
  }

  private setData(data: T): void {
    this.data = data;
    this.notifyCallbacks();
  }

  private setLoading(loading: boolean): void {
    this.loading = loading;
    this.notifyCallbacks();
  }

  private setError(error: string | null): void {
    this.error = error;
    this.notifyCallbacks();
  }

  private notifyCallbacks(): void {
    this.callbacks.forEach(callback => callback());
  }

  subscribe(callback: () => void): () => void {
    this.callbacks.add(callback);
    return () => this.callbacks.delete(callback);
  }

  getData(): T | null { return this.data; }
  isLoading(): boolean { return this.loading; }
  getError(): string | null { return this.error; }
}