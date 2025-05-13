import React, { useState } from 'react';
import { useCreateModel } from '../hooks/api/useAIModels'; // <- You must define this hook
import { QueryClient } from '@tanstack/react-query';

interface CreateModelModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export default function CreateModelModal({ isOpen, onClose }: CreateModelModalProps) {
  const [formData, setFormData] = useState({
    models_name: '',
    number_requests_used: 1000,
    percent_train_data: 80.0,
    percent_normal_requests: 70.0,
    num_trees: 100,
    max_depth: 10,
    max_features: 'sqrt',
    min_samples_split: 5,
    min_samples_leaf: 2 ,
    criterion: 'gini',
  });

  const [errors, setErrors] = useState<Record<string, string>>({});
  const { mutate } = useCreateModel(); // Custom hook
  const queryClient = new QueryClient();

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: isNaN(Number(value)) ? value : Number(value),
    }));
  };

  const validate = () => {
    const newErrors: Record<string, string> = {};
    if (formData.number_requests_used < 1) newErrors.number_requests_used = 'Must be at least 1';
    if (formData.percent_train_data <= 0 || formData.percent_train_data >= 100)
      newErrors.percentTrainData = 'Must be between 1 and 99';
    if (formData.percent_normal_requests < 0 || formData.percent_normal_requests > 100)
      newErrors.percent_normal_requests = 'Must be between 0 and 100';
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate()) return;

    mutate(formData, {
      onSuccess: () => {
        alert('Model created successfully!');
        queryClient.invalidateQueries({ queryKey: ['aiModels'] });
        onClose();
      },
      onError: (error: Error) => {
        alert('Error creating model: ' + error.message);
      },
    });
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div className="bg-white p-6 rounded-lg shadow-lg w-full max-w-3xl">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold">Create AI Model</h2>
          <button onClick={onClose} className="text-gray-500 hover:text-gray-700">X</button>
        </div>

        <form onSubmit={handleSubmit}>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {Object.entries({
              models_name: 'Model Name',
              number_requests_used: 'Number of Requests Used',
              percent_train_data: 'Percent Train Data',
              percent_normal_requests: 'Percent Normal Requests',
              num_trees: 'Number of Trees',
              max_depth: 'Max Depth',
              min_samples_split: 'Min Samples Split',
              min_samples_leaf: 'Min Samples Leaf',
            }).map(([key, label]) => (
              <div key={key}>
                <label className="block text-sm font-medium text-gray-700 mb-1">{label}</label>
                <input
                  name={key}
                  value={(formData as any)[key]}
                  onChange={handleChange}
                  className="w-full px-4 py-2 border border-gray-300 rounded"
                />
                {errors[key] && <p className="text-red-500 text-sm mt-1">{errors[key]}</p>}
              </div>
            ))}

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Max Features</label>
              <select
                name="maxFeatures"
                value={formData.max_features}
                onChange={handleChange}
                className="w-full px-4 py-2 border border-gray-300 rounded"
              >
                <option value="sqrt">sqrt</option>
                <option value="log2">log2</option>
                <option value="none">none</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Criterion</label>
              <select
                name="criterion"
                value={formData.criterion}
                onChange={handleChange}
                className="w-full px-4 py-2 border border-gray-300 rounded"
              >
                <option value="gini">gini</option>
                <option value="entropy">entropy</option>
              </select>
            </div>
          </div>

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
              Create Model
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
