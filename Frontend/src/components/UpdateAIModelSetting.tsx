import React, { useEffect, useState } from 'react';
import { useUpdateAiModelSetting } from '../hooks/api/useAIModels';
import { QueryClient } from '@tanstack/react-query';

interface CreateModelModalProps {
  isOpen: boolean;
  onClose: () => void;
  aiModelSetting: AIModel;
}

interface AIModel {
  id: string;
  expected_accuracy: number;
  expected_precision: number;
  expected_recall: number;
  expected_f1: number;
  train_every: number;
}

export default function UpdateAIModelSetting({ isOpen, onClose, aiModelSetting }: CreateModelModalProps) {
  const [formData, setFormData] = useState({
    id: '',
    expected_accuracy: 0,
    expected_precision: 0,
    expected_recall: 0,
    expected_f1: 0,
    train_every: 0,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});
  const { mutate } = useUpdateAiModelSetting(); // Custom hook
  const queryClient = new QueryClient();

  useEffect(() => {
    if (aiModelSetting) {
      setFormData({
        id: aiModelSetting.id,
        expected_accuracy: aiModelSetting.expected_accuracy || 0,
        expected_precision: aiModelSetting.expected_precision || 0,
        expected_recall: aiModelSetting.expected_recall || 0,
        expected_f1: aiModelSetting.expected_f1 || 0,
        train_every: aiModelSetting.train_every || 0,
      });
    }
  }, [aiModelSetting]);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: parseFloat(value),
    }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    // Optionally add form validation here
    mutate(formData, {
      onSuccess: () => {
        alert('Model settings updated successfully!');
        queryClient.invalidateQueries({ queryKey: ['aiModels'] });
        onClose();
      },
      onError: (error: Error) => {
        alert('Error updating model settings: ' + error.message);
      },
    });
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div className="bg-white p-6 rounded-lg shadow-lg w-full max-w-lg">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold">Update AI Model Settings</h2>
          <button onClick={onClose} className="text-gray-500 hover:text-gray-700">X</button>
        </div>

        <form onSubmit={handleSubmit}>
          {[
            { key: 'expected_accuracy', label: 'Expected Accuracy' },
            { key: 'expected_precision', label: 'Expected Precision' },
            { key: 'expected_recall', label: 'Expected Recall' },
            { key: 'expected_f1', label: 'Expected F1' },
            { key: 'train_every', label: 'Train Every (ms)' },
          ].map(({ key, label }) => (
            <div className="mb-4" key={key}>
              <label className="block text-sm font-medium text-gray-700 mb-1">{label}</label>
              <input
                type="number"
                name={key}
                value={(formData as any)[key]}
                onChange={handleChange}
                className="w-full px-4 py-2 border border-gray-300 rounded"
                step="any"
              />
              {errors[key] && <p className="text-red-500 text-sm mt-1">{errors[key]}</p>}
            </div>
          ))}

          <div className="mt-6 flex justify-end space-x-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              Update Settings
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
