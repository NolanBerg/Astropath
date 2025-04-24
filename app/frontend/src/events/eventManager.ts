type EventCallback<T> = (data: T) => void;

// Event handler class handles incoming events from Golang.
class EventManager {
  private eventHandlers: Map<string, Set<EventCallback<any>>> = new Map();

  constructor() {
    // Initialize all event listeners from runtime
    window["runtime"].EventsOn("*", (event: string, data: any) => {
      const handlers = this.eventHandlers.get(event);
      if (handlers) {
        handlers.forEach((handler) => handler(data));
      }
    });
  }

  public addListener<T>(event: string, callback: EventCallback<T>): void {
    if (!this.eventHandlers.has(event)) {
      this.eventHandlers.set(event, new Set());
    }
    this.eventHandlers.get(event)!.add(callback);
  }

  public removeListener<T>(event: string, callback: EventCallback<T>): void {
    const handlers = this.eventHandlers.get(event);
    if (handlers) {
      handlers.delete(callback);
      if (handlers.size === 0) {
        this.eventHandlers.delete(event);
      }
    }
  }

  public emit<T>(event: string, data: T): void {
    window["runtime"].EventsEmit(event, data);
  }
}

export const eventManager = new EventManager();
