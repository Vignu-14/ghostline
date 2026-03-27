import { Component, type ErrorInfo, type PropsWithChildren } from "react";

type ErrorBoundaryState = {
  hasError: boolean;
  message: string;
};

export class ErrorBoundary extends Component<PropsWithChildren, ErrorBoundaryState> {
  state: ErrorBoundaryState = {
    hasError: false,
    message: "",
  };

  static getDerivedStateFromError(error: Error) {
    return {
      hasError: true,
      message: error.message || "Unexpected frontend error.",
    };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error("Frontend error boundary caught an error:", error, info);
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="center-stage">
          <div className="panel">
            <p className="eyebrow">Frontend error</p>
            <h2>Something slipped out of the feed.</h2>
            <p>{this.state.message || "Refresh the page and try again."}</p>
            <button className="button button--primary" onClick={() => window.location.reload()} type="button">
              Reload app
            </button>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}
