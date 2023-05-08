""" Setup """
import setuptools


def long_description():
  """ Return long description """
  with open('README.md', 'r', encoding='utf-8') as fh:
    return fh.read()


setuptools.setup(
  name="layrz-forms",
  version="1.0.6",
  author="Golden M",
  author_email="software@goldenmcorp.com",
  description="Layrz forms and tools for Python",
  long_description=long_description(),
  long_description_content_type="text/markdown",
  packages=setuptools.find_packages(),
  classifiers=[
    "Programming Language :: Python :: 3",
    "Programming Language :: Python :: 3.8",
    "Programming Language :: Python :: 3.9",
    "Programming Language :: Python :: 3.10",
    "Programming Language :: Python :: 3.11",
  ],
  python_requires=">=3.8",
)
