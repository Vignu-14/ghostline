import type { InputHTMLAttributes } from "react";

type InputProps = InputHTMLAttributes<HTMLInputElement> & {
  label: string;
  error?: string;
};

export function Input({ label, error, className = "", ...props }: InputProps) {
  return (
    <label className="field">
      <span className="field__label">{label}</span>
      <input className={`field__input ${className}`.trim()} {...props} />
      {error ? <span className="field__error">{error}</span> : null}
    </label>
  );
}
