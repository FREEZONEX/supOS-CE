export const getNextStep = (currentStep: string | number, data: any) => {
  const nextStep = data?.steps?.find((f: any) => f.stepName === currentStep)?.nextStep;
  if (!nextStep) {
    return null;
  } else {
    return {
      ...data,
      currentStep: nextStep,
    };
  }
};
